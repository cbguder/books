package goodreads

import (
	"context"
	"net/http"

	"github.com/cbguder/books/api_client"
	"github.com/cbguder/books/config"
)

const (
	baseUrl   = "https://www.goodreads.com"
	userAgent = "Goodreads/4.9.0 (iPhone; iOS 16.6.1; Scale/3.00)"
)

type Client struct {
	apiClient *api_client.ApiClient

	accessToken  string
	refreshToken string
	userId       string
}

func NewClient() *Client {
	cfg := config.Get()

	return &Client{
		apiClient:    api_client.NewApiClient(),
		accessToken:  cfg.Goodreads.AccessToken,
		refreshToken: cfg.Goodreads.RefreshToken,
		userId:       cfg.Goodreads.UserId,
	}
}

func (c *Client) request(ctx context.Context, method, url string, body any) (*http.Request, error) {
	req, err := c.apiClient.Request(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-amz-access-token", c.accessToken)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-APPLE-DEVICE-MODEL", "iPhone")
	req.Header.Set("X_APPLE_SYSTEM_VERSION", "16.6.1")
	req.Header.Set("X_APPLE_APP_VERSION", "900")

	return req, nil
}
