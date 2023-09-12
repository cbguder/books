package goodreads

import (
	"context"
	"net/url"
)

func (c *Client) UpdateUserStatus(ctx context.Context, bookId, page, percent, body string) error {
	val := url.Values{}
	val.Set("_nc", "true")
	val.Set("format", "xml")
	val.Set("user_status[book_id]", bookId)

	if body != "" {
		val.Set("user_status[body]", body)
	}

	if page != "" {
		val.Set("user_status[page]", page)
	}

	if percent != "" {
		val.Set("user_status[percent]", percent)
	}

	loc := "https://www.goodreads.com/user_status.xml"

	req, err := c.request(ctx, "POST", loc, val)
	if err != nil {
		return err
	}

	return c.apiClient.Do(req, nil)
}
