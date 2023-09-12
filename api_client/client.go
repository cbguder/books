package api_client

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ApiClient struct {
	client        *http.Client
	authorization string
}

type Decoder interface {
	Decode(v any) error
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if response == nil {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")

	var decoder Decoder
	if strings.HasPrefix(contentType, "application/json") {
		decoder = json.NewDecoder(resp.Body)
	} else if strings.HasPrefix(contentType, "application/xml") {
		decoder = xml.NewDecoder(resp.Body)
	}

	if decoder == nil {
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	return decoder.Decode(response)
}

func (c *ApiClient) Request(ctx context.Context, method, url string, body any) (*http.Request, error) {
	bodyReader, contentType, err := readerForBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if c.authorization != "" {
		req.Header.Set("Authorization", c.authorization)
	}

	return req, nil
}

func readerForBody(body any) (io.Reader, string, error) {
	var contentType string
	var reader io.Reader

	if body == nil {
		return nil, "", nil
	}

	switch v := body.(type) {
	case io.Reader:
		return v, "", nil

	case url.Values:
		contentType = "application/x-www-form-urlencoded"
		reader = strings.NewReader(v.Encode())

	default:
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, "", err
		}

		contentType = "application/json"
		reader = bytes.NewReader(bodyBytes)
	}

	return reader, contentType, nil
}
