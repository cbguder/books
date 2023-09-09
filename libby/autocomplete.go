package libby

import (
	"context"
	"net/url"
)

type AutocompleteResponse struct {
	Branches []Branch `json:"branches"`
}

type Branch struct {
	Systems []System `json:"systems"`
}

type System struct {
	WebsiteId int `json:"websiteId"`
}

func (c *Client) Autocomplete(ctx context.Context, query string) (*AutocompleteResponse, error) {
	loc := baseUrl + "/locate/autocomplete/" + url.PathEscape(query)

	req, err := c.apiClient.Request(ctx, "GET", loc, nil)
	if err != nil {
		return nil, err
	}

	resp := AutocompleteResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}
