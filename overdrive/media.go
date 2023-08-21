package overdrive

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MediaResponse struct {
	Items []MediaItem `json:"items"`
}

type MediaItem struct {
	EstimatedWaitDays int        `json:"estimatedWaitDays"`
	FirstCreatorName  string     `json:"firstCreatorName"`
	Id                string     `json:"id"`
	IsAvailable       bool       `json:"isAvailable"`
	Languages         []Language `json:"languages"`
	PublishDate       time.Time  `json:"publishDate"`
	Title             string     `json:"title"`
	Type              MediaType  `json:"type"`
}

type Language struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type MediaType struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (c *Client) GetMedia(ctx context.Context, library, query string) (*MediaResponse, error) {
	req, err := mediaQueryRequest(ctx, library, query)
	if err != nil {
		return nil, err
	}

	resp := MediaResponse{}
	err = c.do(req, &resp)
	return &resp, err
}

func mediaQueryRequest(ctx context.Context, library, query string) (*http.Request, error) {
	vls := url.Values{}
	vls.Set("query", query)
	vls.Set("format", "ebook-overdrive,ebook-media-do,ebook-overdrive-provisional,audiobook-overdrive,audiobook-overdrive-provisional,magazine-overdrive")
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
