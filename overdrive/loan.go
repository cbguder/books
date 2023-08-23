package overdrive

import (
	"context"
	"fmt"
	"net/http"
)

type CreateLoanRequest struct {
	Period           int              `json:"period"`
	Units            string           `json:"units"`
	LuckyDay         interface{}      `json:"luckyDay"`
	TitleFormat      string           `json:"title_format"`
	ReportingContext ReportingContext `json:"reporting_context"`
}

type ReportingContext struct {
	ListSourceName string `json:"listSourceName"`
	ListSourceId   string `json:"listSourceId"`
	ListPath       string `json:"listPath"`
	ClientName     string `json:"clientName"`
	ClientVersion  string `json:"clientVersion"`
	Environment    string `json:"environment"`
}

type OpenBookResponse struct {
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

	return c.do(req, nil)
}

func (c *Client) DeleteLoan(ctx context.Context, cardId, mediaId string) error {
	loc := fmt.Sprintf("%s/card/%s/loan/%s", sentry, cardId, mediaId)
	req, err := c.request(ctx, "DELETE", loc, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) OpenBook(ctx context.Context, cardId, mediaId, mediaTypeId string) (*OpenBookResponse, error) {
	req, err := c.openBookRequest(ctx, cardId, mediaId, mediaTypeId)
	if err != nil {
		return nil, err
	}

	resp := OpenBookResponse{}
	err = c.do(req, &resp)
	return &resp, err
}

func (c *Client) GetRosters(ctx context.Context, rostersUrl, message string) ([]Roster, error) {
	rostersUrl = rostersUrl + "?" + message

	req, err := http.NewRequestWithContext(ctx, "GET", rostersUrl, nil)
	if err != nil {
		return nil, err
	}

	var rosters []Roster
	err = c.do(req, &rosters)
	return rosters, err
}

func (c *Client) Download(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, err
}

func (c *Client) createLoanRequest(ctx context.Context, cardId, mediaId, format string) (*http.Request, error) {
	body := CreateLoanRequest{
		Period:      21,
		Units:       "day",
		LuckyDay:    nil,
		TitleFormat: format,
		ReportingContext: ReportingContext{
			ListSourceName: "similar",
			ListSourceId:   mediaId,
			ListPath:       "library/sfpl/similar-" + mediaId,
			ClientName:     "Dewey",
			ClientVersion:  "15.2.1",
			Environment:    "charlie",
		},
	}

	loc := fmt.Sprintf("%s/card/%s/loan/%s", sentry, cardId, mediaId)

	return c.request(ctx, "POST", loc, body)
}

func (c *Client) openBookRequest(ctx context.Context, cardId string, mediaId string, mediaTypeId string) (*http.Request, error) {
	var loc string
	if mediaTypeId == "ebook" {
		loc = fmt.Sprintf("%s/open/book/card/%s/title/%s", sentry, cardId, mediaId)
	} else {
		loc = fmt.Sprintf("%s/open/%s/card/%s/title/%s", sentry, mediaTypeId, cardId, mediaId)
	}

	return c.request(ctx, "GET", loc, nil)
}
