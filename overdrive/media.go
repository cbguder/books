package overdrive

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

type GetMediaResponse struct {
	Id    string    `json:"id"`
	Title string    `json:"title"`
	Type  MediaType `json:"type"`

	EstimatedWaitDays int          `json:"estimatedWaitDays"`
	FirstCreatorName  string       `json:"firstCreatorName"`
	IsAvailable       bool         `json:"isAvailable"`
	PublishDate       FlexibleTime `json:"publishDate"`

	Languages []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"languages"`
}

type SearchMediaResponse struct {
	Items []GetMediaResponse `json:"items"`
}

type MediaType struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (c *Client) GetMedia(ctx context.Context, library, titleId string) (*GetMediaResponse, error) {
	req, err := getMediaRequest(ctx, library, titleId)
	if err != nil {
		return nil, err
	}

	resp := GetMediaResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}

func (c *Client) SearchMedia(ctx context.Context, library, query string, format MediaFormat) (*SearchMediaResponse, error) {
	req, err := searchMediaRequest(ctx, library, query, format)
	if err != nil {
		return nil, err
	}

	resp := SearchMediaResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}

func getMediaRequest(ctx context.Context, library, titleId string) (*http.Request, error) {
	loc := fmt.Sprintf("%s/libraries/%s/media/%s", thunder, library, titleId)
	return http.NewRequestWithContext(ctx, "GET", loc, nil)
}

func searchMediaRequest(ctx context.Context, library, query string, format MediaFormat) (*http.Request, error) {
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
