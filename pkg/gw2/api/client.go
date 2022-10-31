package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	token      string
	httpClient *http.Client
}

func NewGW2APIClient(token string) *Client {
	t := http.DefaultTransport.(*http.Transport).Clone()

	return &Client{
		token: token,
		httpClient: &http.Client{
			Transport: t,
			Timeout:   10 * time.Second,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method string, uri string, object interface{}) error {
	u, err := url.JoinPath("https://api.guildwars2.com", uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Schema-Version", "2022-03-23T19:00:00.000Z")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		apiError, err := NewError(resp)
		if err != nil {
			return err
		}

		return apiError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}
