package api

import (
	"context"
	"net/http"
	"time"
)

type UserDataRecordResponse struct {
	LastUpdated time.Time `json:"LastUpdated"`
	Permission  string    `json:"Permission"`
	Value       string    `json:"Value"`
}

type UpdateUserDataRequest struct {
	CustomTags   map[string]string `json:"CustomTags"`
	Data         map[string]string `json:"Data"`
	KeysToRemove []string          `json:"KeysToRemove"`
	Permission   string            `json:"Permission"`
}

type UpdateUserDataResponse struct {
	DataVersion int `json:"DataVersion"`
}

func (c *Client) UpdateUserData(ctx context.Context, request *UpdateUserDataRequest) (*UpdateUserDataResponse, error) {
	response := &UpdateUserDataResponse{}

	err := c.doRequest(ctx, http.MethodPost, "/Client/UpdateUserData", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
