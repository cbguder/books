package overdrive

import (
	"context"
	"time"
)

type ChipResponse struct {
	Chip     string `json:"chip"`
	Identity string `json:"identity"`
	Syncable bool   `json:"syncable"`
	Primary  bool   `json:"primary"`
}

type CloneRequest struct {
	Code string `json:"code"`
}

type CloneResponse struct {
	Result string `json:"result"`
	Chip   string `json:"chip"`
}

type SyncResponse struct {
	Result   string `json:"result"`
	Cards    []Card `json:"cards"`
	Loans    []Loan `json:"loans"`
	Holds    []Hold `json:"holds"`
	Identity string `json:"identity"`
}

type Card struct {
	CardId       string `json:"cardId"`
	CardName     string `json:"cardName"`
	AdvantageKey string `json:"advantageKey"`

	Library struct {
		Name      string `json:"name"`
		WebsiteId string `json:"websiteId"`
		Logo      struct {
			Url string `json:"url"`
		} `json:"logo"`
	} `json:"library"`
}

type Loan struct {
	CardId           string    `json:"cardId"`
	CheckoutDate     time.Time `json:"checkoutDate"`
	ExpireDate       time.Time `json:"expireDate"`
	FirstCreatorName string    `json:"firstCreatorName"`
	Id               string    `json:"id"`
	Title            string    `json:"title"`
	Type             MediaType `json:"type"`
}

type Hold struct {
	Id                string    `json:"id"`
	CardId            string    `json:"cardId"`
	EstimatedWaitDays int       `json:"estimatedWaitDays"`
	FirstCreatorName  string    `json:"firstCreatorName"`
	PlacedDate        time.Time `json:"placedDate"`
	Title             string    `json:"title"`
	Type              MediaType `json:"type"`
}

func (c *Client) Chip(ctx context.Context) (*ChipResponse, error) {
	req, err := c.apiClient.Request(ctx, "POST", sentry+"/chip", nil)
	if err != nil {
		return nil, err
	}

	resp := ChipResponse{}
	err = c.apiClient.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	c.apiClient.SetAuthorization("Bearer " + resp.Identity)

	return &resp, nil
}

func (c *Client) ChipClone(ctx context.Context, code string) (*CloneResponse, error) {
	reqBody := CloneRequest{
		Code: code,
	}

	req, err := c.apiClient.Request(ctx, "POST", sentry+"/chip/clone/code", reqBody)
	if err != nil {
		return nil, err
	}

	resp := CloneResponse{}
	err = c.apiClient.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) ChipSync(ctx context.Context) (*SyncResponse, error) {
	req, err := c.apiClient.Request(ctx, "GET", sentry+"/chip/sync", nil)
	if err != nil {
		return nil, err
	}

	resp := SyncResponse{}
	err = c.apiClient.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	c.apiClient.SetAuthorization("Bearer " + resp.Identity)

	return &resp, nil
}
