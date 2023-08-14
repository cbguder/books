package overdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func NewClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}

type Client struct {
	client *http.Client
}

func (c *Client) GetMedia(ctx context.Context, library, query string) (*MediaResponse, error) {
	req, err := mediaQueryRequest(ctx, library, query)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	mediaResponse := MediaResponse{}
	json.NewDecoder(resp.Body).Decode(&mediaResponse)

	return &mediaResponse, nil
}

func mediaQueryRequest(ctx context.Context, library, query string) (*http.Request, error) {
	vls := url.Values{}
	vls.Set("query", query)
	vls.Set("format", "ebook-overdrive,ebook-media-do,ebook-overdrive-provisional,audiobook-overdrive,audiobook-overdrive-provisional,magazine-overdrive")
	vls.Set("perPage", "24")
	vls.Set("page", "1")
	vls.Set("x-client-id", "dewey")

	baseUrl := fmt.Sprintf("https://thunder.api.overdrive.com/v2/libraries/%s/media", library)

	req, err := http.NewRequestWithContext(ctx, "GET", baseUrl, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = vls.Encode()

	return req, nil
}
