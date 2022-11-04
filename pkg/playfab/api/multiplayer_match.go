package api

import (
	"context"
	"net/http"
)

type ServerPort struct {
	Name     string `json:"Name"`
	Num      int    `json:"Num"`
	Protocol string `json:"Protocol"`
}

type ServerDetails struct {
	FQDN        string       `json:"Fqdn"`
	IPV4Address string       `json:"IPV4Address"`
	Ports       []ServerPort `json:"Ports"`
	Region      string       `json:"Region"`
}

type GetMatchRequest struct {
	EscapeObject           bool              `json:"EscapeObject"`
	MatchId                string            `json:"MatchId"`
	QueueName              string            `json:"QueueName"`
	ReturnMemberAttributes bool              `json:"ReturnMemberAttributes"`
	CustomTags             map[string]string `json:"CustomTags"`
}

type MatchmakingPlayerWithTeamAssignment struct {
	Attributes MatchmakingPlayerAttributes `json:"Attributes"`
	Entity     EntityKey                   `json:"Entity"`
	TeamId     *string                     `json:"TeamId"`
}

type GetMatchResponse struct {
	ArrangementString string         `json:"ArrangementString"`
	MatchId           string         `json:"MatchId"`
	Members           []interface{}  `json:"Members"`
	RegionPreferences []string       `json:"RegionPreferences"`
	ServerDetails     *ServerDetails `json:"ServerDetails"`
}

func (c *Client) GetMatch(ctx context.Context, entityToken EntityToken, request *GetMatchRequest) (*GetMatchResponse, error) {
	response := &GetMatchResponse{}

	header := make(http.Header)
	header.Add("X-EntityToken", entityToken.EntityToken)

	err := c.doRequest(ctx, http.MethodPost, "/Match/GetMatch", header, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
