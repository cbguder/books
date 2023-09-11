package goodreads

import (
	"context"
	"encoding/xml"
	"net/url"
)

type CurrentUserDataResponse struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`

	CurrentUser struct {
		User struct {
			ID string `xml:"id"`
		} `xml:"user"`

		CurrentlyReading []struct {
			Title   string `xml:"title"`
			Authors []struct {
				Name string `xml:"name"`
			} `xml:"authors>author"`
		} `xml:"currently_reading>book"`
	} `xml:"current_user"`
}

func (c *Client) GetCurrentUserData(ctx context.Context) (*CurrentUserDataResponse, error) {
	val := url.Values{}
	val.Set("_nc", "true")
	val.Set("format", "xml")
	val.Set("v", "2")

	loc := baseUrl + "/api/current_user_data?" + val.Encode()
	req, err := c.request(ctx, "GET", loc, nil)
	if err != nil {
		return nil, err
	}

	resp := CurrentUserDataResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}
