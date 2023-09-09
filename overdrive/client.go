package overdrive

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/cbguder/books/api_client"
)

const thunder = "https://thunder.api.overdrive.com/v2"
const sentry = "https://sentry-read.svc.overdrive.com"

func NewClient(identity string) *Client {
	jar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar: jar,
	}

	apiClient := api_client.NewApiClientWithHttpClient(httpClient)
	apiClient.SetAuthorization("Bearer " + identity)

	return &Client{
		identity:   identity,
		httpClient: httpClient,
		apiClient:  apiClient,
	}
}

type Client struct {
	identity string

	httpClient *http.Client
	apiClient  *api_client.ApiClient
}
