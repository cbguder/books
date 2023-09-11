package goodreads

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
)

type ReviewsResponse struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`

	Reviews []struct {
		ID   string `xml:"id"`
		Book struct {
			Title           string `xml:"title"`
			PublicationYear string `xml:"publication_year"`
			Authors         []struct {
				Name string `xml:"name"`
			} `xml:"authors>author"`
		} `xml:"book"`
	} `xml:"reviews>review"`
}

func (c *Client) GetReviews(ctx context.Context, userId, shelfName string) (*ReviewsResponse, error) {
	val := url.Values{}
	val.Set("_nc", "true")
	val.Set("format", "xml")
	val.Set("order", "d")
	val.Set("page", "1")
	val.Set("per_page", "100")
	val.Set("shelf", shelfName)
	val.Set("sort", "date_updated")
	val.Set("v", "2")

	loc := fmt.Sprintf("%s/review/list/%s?%s", baseUrl, userId, val.Encode())
	req, err := c.request(ctx, "GET", loc, nil)
	if err != nil {
		return nil, err
	}

	resp := ReviewsResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}
