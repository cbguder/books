package libby

import "github.com/cbguder/books/api_client"

const baseUrl = "https://libbyapp.com/api"

func NewClient() *Client {
	return &Client{
		apiClient: api_client.NewApiClient(),
	}
}

type Client struct {
	apiClient *api_client.ApiClient
}
