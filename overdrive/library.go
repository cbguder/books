package overdrive

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type LibrariesResponse struct {
	Items []struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		PreferredKey string `json:"preferredKey"`
	} `json:"items"`
}

func (c *Client) GetLibrariesByWebsiteId(ctx context.Context, websiteIds []int) (*LibrariesResponse, error) {
	req, err := librariesByWebsiteIdRequest(ctx, websiteIds)
	if err != nil {
		return nil, err
	}

	resp := LibrariesResponse{}
	err = c.do(req, &resp)
	return &resp, err
}

func librariesByWebsiteIdRequest(ctx context.Context, websiteIds []int) (*http.Request, error) {
	websiteIdStrs := make([]string, len(websiteIds))
	for i, websiteId := range websiteIds {
		websiteIdStrs[i] = strconv.Itoa(websiteId)
	}

	websiteIdsStr := strings.Join(websiteIdStrs, ",")

	loc := thunder + "/libraries?websiteIds=" + websiteIdsStr
	return http.NewRequestWithContext(ctx, "GET", loc, nil)
}
