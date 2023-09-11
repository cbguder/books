package goodreads

import (
	"context"
	"encoding/xml"
	"net/url"
)

type CurrentUserShelvesResponse struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`

	CurrentUser struct {
		UserShelves []struct {
			Name          string `xml:"name"`
			BookCount     string `xml:"book_count"`
			ExclusiveFlag string `xml:"exclusive_flag"`
		} `xml:"user_shelves>user_shelf"`
	} `xml:"current_user"`
}

func (c *Client) GetCurrentUserShelves(ctx context.Context) (*CurrentUserShelvesResponse, error) {
	val := url.Values{}
	val.Set("_nc", "true")
	val.Set("format", "xml")

	loc := baseUrl + "/api/current_user_shelves?" + val.Encode()
	req, err := c.request(ctx, "GET", loc, nil)
	if err != nil {
		return nil, err
	}

	resp := CurrentUserShelvesResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}
