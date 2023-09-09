package api_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiClient struct {
	client        *http.Client
	authorization string
}

func NewApiClient() *ApiClient {
	return NewApiClientWithHttpClient(http.DefaultClient)
}

func NewApiClientWithHttpClient(client *http.Client) *ApiClient {
	return &ApiClient{
		client: client,
	}
}

func (c *ApiClient) SetAuthorization(authorization string) {
	c.authorization = authorization
}

func (c *ApiClient) Do(req *http.Request, response any) error {
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

func (c *ApiClient) Request(ctx context.Context, method, url string, body any) (*http.Request, error) {
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

	if c.authorization != "" {
		req.Header.Set("Authorization", c.authorization)
	}

	return req, nil
}
