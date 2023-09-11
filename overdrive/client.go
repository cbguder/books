package overdrive

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/cbguder/books/api_client"
	"github.com/cbguder/books/config"
)

const thunder = "https://thunder.api.overdrive.com/v2"
const sentry = "https://sentry-read.svc.overdrive.com"

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar: jar,
	}

	apiClient := api_client.NewApiClientWithHttpClient(httpClient)

	cfg := config.Get()
	if cfg.Identity != "" {
		apiClient.SetAuthorization("Bearer " + cfg.Identity)
	}

	return &Client{
		httpClient: httpClient,
		apiClient:  apiClient,
	}
}

type Client struct {
	httpClient *http.Client
	apiClient  *api_client.ApiClient
}
