package playfab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	titleId    string
	httpClient *http.Client
}

func NewPlayFabClient() *Client {
	t := http.DefaultTransport.(*http.Transport).Clone()

	return &Client{
		titleId: os.Getenv("PLAYFAB_TITLE_ID"),
		httpClient: &http.Client{
			Transport: t,
			Timeout:   10 * time.Second,
		},
	}
}

type APIResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func (c *Client) doRequest(ctx context.Context, method string, uri string, header http.Header, reqBody interface{}, object interface{}) error {
	u, err := url.JoinPath(fmt.Sprintf("https://%s.playfabapi.com", c.titleId), uri)
	if err != nil {
		return err
	}

	var bodyReader io.Reader

	if reqBody != nil {
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("error marshaling body: %w", err)
		}

		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header = header
	if reqBody != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		errorResponse := &APIErrorResponse{}
		err = json.Unmarshal(body, errorResponse)
		if err != nil {
			return fmt.Errorf("error unmarshaling error: %w, http code: %d, body: %s", err, resp.StatusCode, string(body))
		}

		return errorResponse
	}

	apiResponse := &APIResponse{}
	err = json.Unmarshal(body, apiResponse)
	if err != nil {
		return fmt.Errorf("error unmarshaling body: %w, body: %s", err, string(body))
	}

	if object != nil {
		rawResponseData, err := json.Marshal(apiResponse.Data)
		if err != nil {
			return fmt.Errorf("error marshaling api response data: %w, body: %s", err, apiResponse.Data)
		}

		err = json.Unmarshal(rawResponseData, object)
		if err != nil {
			return fmt.Errorf("error unmarshaling api response data: %w, body: %s", err, string(rawResponseData))
		}
	}

	return nil
}
