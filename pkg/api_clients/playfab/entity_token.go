package playfab

import (
	"context"
	"net/http"
	"os"
)

type EntityKey struct {
	Id   string `json:"Id"`
	Type string `json:"Type"`
}

type EntityToken struct {
	Entity          EntityKey `json:"Entity"`
	EntityToken     string    `json:"EntityToken"`
	TokenExpiration string    `json:"TokenExpiration"`
}

type GetEntityTokenRequest struct {
	CustomTags map[string]string `json:"CustomTags"`
	Entity     EntityKey         `json:"Entity"`
}

func (c *Client) GetTitleEntityToken(ctx context.Context) (*EntityToken, error) {
	request := &EntityKey{
		Id:   os.Getenv("PLAYFAB_TITLE_ID"),
		Type: "title",
	}
	response := &EntityToken{}

	header := make(http.Header)
	header.Add("X-SecretKey", os.Getenv("PLAYFAB_TITLE_SECRET_KEY"))

	err := c.doRequest(ctx, http.MethodPost, "/Authentication/GetEntityToken", header, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
