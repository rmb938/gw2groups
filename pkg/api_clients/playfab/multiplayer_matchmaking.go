package playfab

import (
	"context"
	"net/http"
	"time"
)

type MatchmakingPlayerAttributes struct {
	DataObject        map[string]interface{} `json:"DataObject,omitempty"`
	EscapedDataObject *string                `json:"EscapedDataObject,omitempty"`
}

type MatchmakingPlayer struct {
	Attributes MatchmakingPlayerAttributes `json:"Attributes"`
	Entity     EntityKey                   `json:"Entity"`
}

type CreateMatchMakingTicketRequest struct {
	Creator            MatchmakingPlayer `json:"Creator"`
	GiveUpAfterSeconds int64             `json:"GiveUpAfterSeconds"`
	QueueName          string            `json:"QueueName"`
	CustomTags         map[string]string `json:"CustomTags"`
	MembersToMatchWith []EntityKey       `json:"MembersToMatchWith"`
}

type CreateMatchMakingTicketResponse struct {
	TicketId string `json:"TicketId"`
}

func (c *Client) CreateMatchMakingTicket(ctx context.Context, entityToken EntityToken, request *CreateMatchMakingTicketRequest) (*CreateMatchMakingTicketResponse, error) {
	response := &CreateMatchMakingTicketResponse{}

	header := make(http.Header)
	header.Add("X-EntityToken", entityToken.EntityToken)

	err := c.doRequest(ctx, http.MethodPost, "/Match/CreateMatchmakingTicket", header, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type ListMatchmakingTicketsForPlayerRequest struct {
	QueueName  string            `json:"QueueName"`
	CustomTags map[string]string `json:"CustomTags"`
	Entity     EntityKey         `json:"Entity"`
}

type ListMatchmakingTicketsForPlayerResponse struct {
	TicketIds []string `json:"TicketIds"`
}

func (c *Client) ListMatchmakingTicketsForPlayer(ctx context.Context, entityToken EntityToken, request *ListMatchmakingTicketsForPlayerRequest) (*ListMatchmakingTicketsForPlayerResponse, error) {
	response := &ListMatchmakingTicketsForPlayerResponse{}

	header := make(http.Header)
	header.Add("X-EntityToken", entityToken.EntityToken)

	err := c.doRequest(ctx, http.MethodPost, "/Match/ListMatchmakingTicketsForPlayer", header, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type CancelAllMatchmakingTicketsForPlayerRequest struct {
	QueueName  string            `json:"QueueName"`
	CustomTags map[string]string `json:"CustomTags"`
	Entity     EntityKey         `json:"Entity"`
}

func (c *Client) CancelAllMatchmakingTicketsForPlayer(ctx context.Context, entityToken EntityToken, request *CancelAllMatchmakingTicketsForPlayerRequest) error {
	header := make(http.Header)
	header.Add("X-EntityToken", entityToken.EntityToken)

	err := c.doRequest(ctx, http.MethodPost, "/Match/CancelAllMatchmakingTicketsForPlayer", header, request, nil)
	if err != nil {
		return err
	}

	return nil
}

type GetMatchMakingTicketRequest struct {
	EscapeObject bool              `json:"EscapeObject"`
	QueueName    string            `json:"QueueName"`
	TicketId     string            `json:"TicketId"`
	CustomTags   map[string]string `json:"CustomTags"`
}

type GetMatchMakingTicketResponse struct {
	CancellationReasonString string              `json:"CancellationReasonString"`
	ChangeNumber             int                 `json:"ChangeNumber"`
	Created                  time.Time           `json:"Created"`
	Creator                  MatchmakingPlayer   `json:"Creator"`
	GiveUpAfterSeconds       int64               `json:"GiveUpAfterSeconds"`
	MatchId                  *string             `json:"MatchId"`
	Members                  []MatchmakingPlayer `json:"Members"`
	MembersToMatchWith       []EntityKey         `json:"MembersToMatchWith"`
	QueueName                string              `json:"QueueName"`
	Status                   string              `json:"Status"`
	TicketId                 string              `json:"TicketId"`
}

func (c *Client) GetMatchmakingTicket(ctx context.Context, entityToken EntityToken, request *GetMatchMakingTicketRequest) (*GetMatchMakingTicketResponse, error) {
	response := &GetMatchMakingTicketResponse{}

	header := make(http.Header)
	header.Add("X-EntityToken", entityToken.EntityToken)

	err := c.doRequest(ctx, http.MethodPost, "/Match/GetMatchmakingTicket", header, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
