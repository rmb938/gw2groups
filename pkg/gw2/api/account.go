package api

import (
	"context"
	"net/http"
	"time"
)

type AccountResponse struct {
	ID                string    `json:"id"`
	Age               int       `json:"age"`
	Name              string    `json:"name"`
	World             int       `json:"world"`
	Guilds            []string  `json:"guilds"`
	GuildLeader       []string  `json:"guild_leader"`
	Created           time.Time `json:"created"`
	Access            []string  `json:"access"`
	Commander         bool      `json:"commander"`
	FractalLevel      *int      `json:"fractal_level"`
	DailyAP           *int      `json:"daily_ap"`
	MonthlyAP         *int      `json:"monthly_ap"`
	WvWRank           *int      `json:"wvw_rank"`
	LastModified      time.Time `json:"last_modified"`
	BuildStorageSlots *int      `json:"build_storage_slots"`
}

func (c *Client) GetAccount(ctx context.Context) (*AccountResponse, error) {
	response := &AccountResponse{}

	err := c.doRequest(ctx, http.MethodGet, "/v2/account", response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
