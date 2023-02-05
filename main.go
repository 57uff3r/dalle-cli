package main

import (
	dalle "dalle_cli/dalle"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
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
				Usage:     "Export all notes and highlights from book with [BOOK_ID]",
				UsageText: "Export all notes and highlights from book with [BOOK_ID]",
				Action:    generateImage,
				ArgsUsage: "ibooks_notes_exporter export BOOK_ID_GOES_HERE",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "description"},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func generateImage(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		log.Fatal("For generating an image, you have to pass description: dallecli generate DESCRIPTION_GOES_HERE")
	}

	dalle_api_key := os.Getenv("DALLE_API_KEY")
	client := dalle.NewClient(dalle_api_key)
	desciprtion := cCtx.Args().Get(0)

	log.Println("Generating image")
	data, err := client.Generate(desciprtion, nil, nil, nil, nil)

	if err != nil {
		log.Println(err)
	}

	imageURL := data[0].URL

	log.Println("Downloading image")

	// Download image
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bar := progressbar.DefaultBytes(
		response.ContentLength,
		"downloading",
	)

	if response.StatusCode != 200 {
		return errors.New("Received non-200 response code")
	}

	id := uuid.New()
	fname := fmt.Sprintf("%s.png", id.String())

	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(io.MultiWriter(file, bar), response.Body)
	if err != nil {
		return err
	}

	return nil
}
