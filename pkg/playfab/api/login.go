package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"k8s.io/utils/pointer"
)

type LoginWithCustomIDRequest struct {
	TitleId               string                           `json:"TitleId"`
	CreateAccount         *bool                            `json:"CreateAccount,omitempty"`
	CustomId              *string                          `json:"CustomId,omitempty"`
	CustomTags            map[string]string                `json:"CustomTags,omitempty"`
	EncryptedRequest      *string                          `json:"EncryptedRequest,omitempty"`
	InfoRequestParameters *PlayerCombinedInfoRequestParams `json:"InfoRequestParameters,omitempty"`
	PlayerSecret          *string                          `json:"PlayerSecret,omitempty"`
}

type EntityKey struct {
	Id   string `json:"Id"`
	Type string `json:"Type"`
}

type EntityTokenResponse struct {
	Entity          EntityKey `json:"Entity"`
	EntityToken     string    `json:"EntityToken"`
	TokenExpiration string    `json:"TokenExpiration"`
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

type LoginResponse struct {
	EntityToken         EntityTokenResponse         `json:"EntityToken"`
	InfoResultPayload   PlayerCombinedInfoResponse  `json:"InfoResultPayload"`
	LastLoginTime       time.Time                   `json:"LastLoginTime"`
	NewlyCreated        bool                        `json:"NewlyCreated"`
	PlayFabId           string                      `json:"PlayFabId"`
	SessionTicket       string                      `json:"SessionTicket"`
	SettingsForUser     UserSettingsResponse        `json:"SettingsForUser"`
	TreatmentAssignment TreatmentAssignmentResponse `json:"TreatmentAssignment"`
}

func (c *Client) LoginWithCustomID(ctx context.Context, request *LoginWithCustomIDRequest) (*LoginResponse, error) {
	response := &LoginResponse{}

	request.TitleId = os.Getenv("PLAYFAB_TITLE_ID")

	err := c.doRequest(ctx, http.MethodPost, "/Client/LoginWithCustomID", request, response)
	if err != nil {
		return nil, err
	}

	c.sessionTicket = pointer.String(response.SessionTicket)
	c.entityToken = pointer.String(response.EntityToken.EntityToken)

	return response, nil
}
