package playfab

import (
	"context"
	"net/http"
	"os"
	"time"
)

type UserDataRecordResponse struct {
	LastUpdated time.Time `json:"LastUpdated"`
	Permission  string    `json:"Permission"`
	Value       string    `json:"Value"`
}

type ServerUpdateUserDataRequest struct {
	PlayFabId    string            `json:"PlayFabId"`
	CustomTags   map[string]string `json:"CustomTags"`
	Data         map[string]string `json:"Data"`
	KeysToRemove []string          `json:"KeysToRemove"`
	Permission   string            `json:"Permission"`
}

type UpdateUserDataResponse struct {
	DataVersion int `json:"DataVersion"`
}

func (c *Client) UpdateUserData(ctx context.Context, playfabID string, request *ServerUpdateUserDataRequest) (*UpdateUserDataResponse, error) {
	response := &UpdateUserDataResponse{}

	header := make(http.Header)
	header.Add("X-SecretKey", os.Getenv("PLAYFAB_TITLE_SECRET_KEY"))

	request.PlayFabId = playfabID

	err := c.doRequest(ctx, http.MethodPost, "/Server/UpdateUserData", header, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
