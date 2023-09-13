package overdrive

import (
	"context"
	"fmt"
)

type CreateHoldRequest struct{}

func (c *Client) CreateHold(ctx context.Context, cardId, mediaId string) error {
	body := CreateHoldRequest{}

	loc := fmt.Sprintf("%s/card/%s/hold/%s", sentry, cardId, mediaId)

	req, err := c.apiClient.Request(ctx, "POST", loc, body)
	if err != nil {
		return err
	}

	return c.apiClient.Do(req, nil)
}

func (c *Client) DeleteHold(ctx context.Context, cardId, mediaId string) error {
	loc := fmt.Sprintf("%s/card/%s/hold/%s", sentry, cardId, mediaId)

	req, err := c.apiClient.Request(ctx, "DELETE", loc, nil)
	if err != nil {
		return err
	}

	return c.apiClient.Do(req, nil)
}
