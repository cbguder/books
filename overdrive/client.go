package overdrive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

const thunder = "https://thunder.api.overdrive.com/v2"
const sentry = "https://sentry-read.svc.overdrive.com"

func NewClient(identity string) *Client {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: jar,
	}

	return &Client{
		identity: identity,
		client:   client,
	}
}

type Client struct {
	identity string
	client   *http.Client
}

func (c *Client) do(req *http.Request, response any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if response == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(response)
}

func (c *Client) request(ctx context.Context, method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.identity != "" {
		req.Header.Set("Authorization", "Bearer "+c.identity)
	}

	return req, nil
}
