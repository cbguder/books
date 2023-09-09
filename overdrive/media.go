package overdrive

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const audiobookFormats = "audiobook-overdrive,audiobook-overdrive-provisional"
const ebookFormats = "ebook-overdrive,ebook-media-do,ebook-overdrive-provisional"
const magazineFormats = "magazine-overdrive"
const allFormats = ebookFormats + "," + audiobookFormats + "," + magazineFormats

type MediaFormat int

const (
	MediaFormatAny MediaFormat = iota
	MediaFormatAudiobook
	MediaFormatEbook
)

func (f MediaFormat) queryValue() string {
	if f == MediaFormatAudiobook {
		return audiobookFormats
	} else if f == MediaFormatEbook {
		return ebookFormats
	}

	return allFormats
}

type MediaResponse struct {
	Items []struct {
		Id    string    `json:"id"`
		Title string    `json:"title"`
		Type  MediaType `json:"type"`

		EstimatedWaitDays int       `json:"estimatedWaitDays"`
		FirstCreatorName  string    `json:"firstCreatorName"`
		IsAvailable       bool      `json:"isAvailable"`
		PublishDate       time.Time `json:"publishDate"`

		Languages []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"languages"`
	} `json:"items"`
}

type MediaType struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (c *Client) GetMedia(ctx context.Context, library, query string, format MediaFormat) (*MediaResponse, error) {
	req, err := mediaQueryRequest(ctx, library, query, format)
	if err != nil {
		return nil, err
	}

	resp := MediaResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}

func mediaQueryRequest(ctx context.Context, library, query string, format MediaFormat) (*http.Request, error) {
	vls := url.Values{}
	vls.Set("query", query)
	vls.Set("format", format.queryValue())
	vls.Set("perPage", "24")
	vls.Set("page", "1")
	vls.Set("x-client-id", "dewey")

	loc := fmt.Sprintf("%s/libraries/%s/media", thunder, library)

	req, err := http.NewRequestWithContext(ctx, "GET", loc, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = vls.Encode()

	return req, nil
}
