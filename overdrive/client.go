package overdrive

import (
	"github.com/cbguder/books/api_client"
	"github.com/cbguder/books/config"
)

const thunder = "https://thunder.api.overdrive.com/v2"
const sentry = "https://sentry-read.svc.overdrive.com"

type Client struct {
	apiClient *api_client.ApiClient
}

func NewClient() *Client {
	apiClient := api_client.NewApiClient()

	cfg := config.Get()
	if cfg.Identity != "" {
		apiClient.SetAuthorization("Bearer " + cfg.Identity)
	}

	return &Client{
		apiClient: apiClient,
	}
}
