package dalle_cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Sizes
const (
	Small  int = 256
	Medium int = 512
	Large  int = 1024
)

const (
	defaultBaseURL   = "https://api.openai.com/v1/images"
	defaultUserAgent = "dalle_cli"
	defaultTimeout   = 30 * time.Second
)

type Response struct {
	Created int64   `json:"created"`
	Data    []Datum `json:"data"`
}

type Datum struct {
	URL string `json:"url"`
}

const (
	URLFormat        = "url"
	Base64JSONFormat = "b64_json"
)

type GenerateRequest struct {
	Prompt         string  `json:"prompt"`
	N              *int    `json:"n,omitempty"`
	Size           *string `json:"size,omitempty"`
	ResponseFormat *string `json:"response_format,omitempty"`
	User           *string `json:"user,omitempty"`
}

type Client interface {
	Generate(prompt string, size *int, n *int, user *string, responseType *string) ([]Datum, error)
}

type client struct {
	baseURL    string
	apiKey     string
	userAgent  string
	httpClient *http.Client
}

func NewClient(apiKey string) Client {
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	c := &client{
		baseURL:    defaultBaseURL,
		apiKey:     apiKey,
		userAgent:  defaultUserAgent,
		httpClient: httpClient,
	}

	return c
}

func pointerizeString(s string) *string {
	return &s
}

// Prompt is the prompt to generate an image from.
//
// Size is the size of the image to generate (Small, Medium, Large).
//
// N is the number of images to generate.
//
// https://beta.openai.com/docs/guides/images/usage
func (c *client) Generate(prompt string, size *int, n *int, user *string, responseType *string) ([]Datum, error) {
	url := c.baseURL + "/generations"

	var sizeStr *string

	if size != nil {
		sizeStr = pointerizeString(fmt.Sprintf("%dx%d", size, size))
	}

	body := GenerateRequest{
		Prompt:         prompt,
		N:              n,
		Size:           sizeStr,
		User:           user,
		ResponseFormat: responseType,
	}

	jsonStr, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// http error codes
	httpErrorCodes := make(map[int]string)
	httpErrorCodes[400] = "bad request"
	httpErrorCodes[401] = "unauthorized"
	httpErrorCodes[403] = "forbidden"
	httpErrorCodes[404] = "not found"
	httpErrorCodes[429] = "too many requests"
	httpErrorCodes[500] = "internal server error"
	httpErrorCodes[502] = "bad gateway"
	httpErrorCodes[503] = "service unavailable"
	httpErrorCodes[504] = "gateway timeout"

	if resp.StatusCode != 200 {
		if errorText, ok := httpErrorCodes[resp.StatusCode]; ok {
			log.Fatal(errorText)
		} else {
			log.Fatal("Unknown error")
		}
	}

	var response Response

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
