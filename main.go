package main

import (
	dalle "dalle_cli/dalle"
	"fmt"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	_, found := os.LookupEnv("DALLE_API_KEY")
	if found == false {
		log.Fatal("DALLE_API_KEY env varilable is not set")
	}

	app := &cli.App{
		Name:    "DALL E CLI tool",
		Usage:   "Generate images with DALL E",
		Authors: []*cli.Author{{Name: "Andrey Korchak", Email: "me@akorchak.software"}},
		Commands: []*cli.Command{
			{
				Name:      "generate",
				HideHelp:  false,
				Usage:     "Generate an image from description",
				UsageText: "Generate an image from description",
				ArgsUsage: "dallecli generate --description \"Image description goes here\" --howmany",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "description",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "howmany",
						Value:    2,
						Required: true,
					},
				},
				Action: generateImage,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func downloadImageWorker(URL string, folderName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(fmt.Sprintf("Can't download image from %s, got non-200 response code", URL))
		return nil
	}
	fname := fmt.Sprintf("%s/%s.png", folderName, uuid.New().String())
	file, err := os.Create(fname)
	if err != nil {
		log.Println(fmt.Sprintf("Can't save image from %s, got error %s", URL, err))
		return nil
	}
	defer file.Close()
	_, err = io.Copy(io.MultiWriter(file), response.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Can't save image from %s, got error %s", URL, err))
		return nil
	}

	return nil
}

func generateImage(cCtx *cli.Context) error {
	var imageDescriptionArg string = cCtx.String("description")
	howmanyImagesArg := int(cCtx.Uint("howmany"))

	if howmanyImagesArg > 20 {
		log.Fatal("You are trying to generate more than 20 images. We limited the max number of images that can" +
			" be generated to 20 in order to save you from accidentally spending too much money.\n")
		return nil
	}

	folderForResults := uuid.New().String()
	if err := os.Mkdir(folderForResults, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Your images will be stored in folder %s", folderForResults))

	dalleApiKey := os.Getenv("DALLE_API_KEY")
	client := dalle.NewClient(dalleApiKey)

	log.Println("Generating images...")
	dalleResponse, err := client.Generate(imageDescriptionArg, dalle.Large, howmanyImagesArg, nil, nil)
	if err != nil {
		log.Println(err)
	}

	bar := progressbar.Default(
		int64(len(dalleResponse)),
		"Downloading images",
	)

	var wg sync.WaitGroup
	for i := 0; i < len(dalleResponse); i++ {
		urlToDownload := dalleResponse[i].URL
		wg.Add(1)
		go func() {
			defer wg.Done()
			downloadImageWorker(urlToDownload, folderForResults)
			bar.Add(1)
		}()

	}

	wg.Wait()
	return nil
}
