package playfab

import (
	"context"
	"net/http"
	"os"
	"time"
)

type ServerLoginWithCustomIDRequest struct {
	CreateAccount         *bool                            `json:"CreateAccount,omitempty"`
	CustomTags            map[string]string                `json:"CustomTags,omitempty"`
	InfoRequestParameters *PlayerCombinedInfoRequestParams `json:"InfoRequestParameters,omitempty"`
	PlayerSecret          *string                          `json:"PlayerSecret,omitempty"`
	ServerCustomId        string                           `json:"ServerCustomId"`
}

type UserSettingsResponse struct {
	GatherDeviceInfo bool `json:"GatherDeviceInfo"`
	GatherFocusInfo  bool `json:"GatherFocusInfo"`
	NeedsAttribution bool `json:"NeedsAttribution"`
}

type Variable struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type TreatmentAssignmentResponse struct {
	Variables []Variable `json:"Variables"`
	Variants  []string   `json:"Variants"`
}

type ServerLoginResponse struct {
	EntityToken         EntityToken                 `json:"EntityToken"`
	InfoResultPayload   PlayerCombinedInfoResponse  `json:"InfoResultPayload"`
	LastLoginTime       time.Time                   `json:"LastLoginTime"`
	NewlyCreated        bool                        `json:"NewlyCreated"`
	PlayFabId           string                      `json:"PlayFabId"`
	SessionTicket       string                      `json:"SessionTicket"`
	SettingsForUser     UserSettingsResponse        `json:"SettingsForUser"`
	TreatmentAssignment TreatmentAssignmentResponse `json:"TreatmentAssignment"`
}

func (c *Client) LoginWithCustomID(ctx context.Context, request *ServerLoginWithCustomIDRequest) (*ServerLoginResponse, error) {
	response := &ServerLoginResponse{}

	header := make(http.Header)
	header.Add("X-SecretKey", os.Getenv("PLAYFAB_TITLE_SECRET_KEY"))

	err := c.doRequest(ctx, http.MethodPost, "/Server/LoginWithServerCustomId", header, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
