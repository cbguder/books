package overdrive

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseUrl = "https://thunder.api.overdrive.com/v2"

func NewClient() *Client {
	return &Client{}
}

type Client struct{}

func (c *Client) do(req *http.Request, response any) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(response)
}
