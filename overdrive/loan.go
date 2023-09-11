package overdrive

import (
	"context"
	"fmt"
	"net/http"
)

type CreateLoanRequest struct {
	Period      int         `json:"period"`
	Units       string      `json:"units"`
	LuckyDay    interface{} `json:"luckyDay"`
	TitleFormat string      `json:"title_format"`
}

type OpenLoanResponse struct {
	Urls struct {
		Web        string `json:"web"`
		Openbook   string `json:"openbook"`
		Rosters    string `json:"rosters"`
		Possession string `json:"possession"`
		Activity   string `json:"activity"`
	} `json:"urls"`
	Message string `json:"message"`
}

type Roster struct {
	Group   string `json:"group"`
	Entries []struct {
		Url   string `json:"url"`
		Bytes int    `json:"bytes"`
	} `json:"entries"`
}

func (c *Client) CreateLoan(ctx context.Context, cardId, mediaId, format string) error {
	req, err := c.createLoanRequest(ctx, cardId, mediaId, format)
	if err != nil {
		return err
	}

	return c.apiClient.Do(req, nil)
}

func (c *Client) DeleteLoan(ctx context.Context, cardId, mediaId string) error {
	loc := fmt.Sprintf("%s/card/%s/loan/%s", sentry, cardId, mediaId)
	req, err := c.apiClient.Request(ctx, "DELETE", loc, nil)
	if err != nil {
		return err
	}

	return c.apiClient.Do(req, nil)
}

func (c *Client) OpenLoan(ctx context.Context, cardId, mediaId, mediaTypeId string) (*OpenLoanResponse, error) {
	req, err := c.openBookRequest(ctx, cardId, mediaId, mediaTypeId)
	if err != nil {
		return nil, err
	}

	resp := OpenLoanResponse{}
	err = c.apiClient.Do(req, &resp)
	return &resp, err
}

func (c *Client) createLoanRequest(ctx context.Context, cardId, mediaId, format string) (*http.Request, error) {
	body := CreateLoanRequest{
		Period:      21,
		Units:       "day",
		LuckyDay:    nil,
		TitleFormat: format,
	}

	loc := fmt.Sprintf("%s/card/%s/loan/%s", sentry, cardId, mediaId)

	return c.apiClient.Request(ctx, "POST", loc, body)
}

func (c *Client) openBookRequest(ctx context.Context, cardId string, mediaId string, mediaTypeId string) (*http.Request, error) {
	var loc string
	if mediaTypeId == "ebook" {
		loc = fmt.Sprintf("%s/open/book/card/%s/title/%s", sentry, cardId, mediaId)
	} else {
		loc = fmt.Sprintf("%s/open/%s/card/%s/title/%s", sentry, mediaTypeId, cardId, mediaId)
	}

	return c.apiClient.Request(ctx, "GET", loc, nil)
}
