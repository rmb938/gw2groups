package gw2

import (
	"encoding/json"
	"io"
	"net/http"
)

type StatusReason string

const (
	StatusReasonUnknown      StatusReason = ""
	StatusReasonUnauthorized StatusReason = "Unauthorized"
	StatusReasonForbidden    StatusReason = "Forbidden"
	StatusReasonNotFound     StatusReason = "NotFound"
)

type Error struct {
	Error *string `json:"error"`
	Text  *string `json:"text"`
}

type Status struct {
	Message string
	Reason  StatusReason
	Code    int32
}

type StatusError struct {
	ErrStatus Status
}

type APIStatus interface {
	Status() Status
}

var _ error = &StatusError{}

// Error implements the Error interface.
func (e *StatusError) Error() string {
	return e.ErrStatus.Message
}

// Status allows access to e's status without having to know the detailed workings
// of StatusError.
func (e *StatusError) Status() Status {
	return e.ErrStatus
}

func NewUnauthorized(message string) *StatusError {
	return &StatusError{
		ErrStatus: Status{
			Message: message,
			Reason:  StatusReasonUnauthorized,
			Code:    http.StatusUnauthorized,
		},
	}
}

func NewForbidden(message string) *StatusError {
	return &StatusError{
		ErrStatus: Status{
			Message: message,
			Reason:  StatusReasonForbidden,
			Code:    http.StatusForbidden,
		},
	}
}

func NewNotFound(message string) *StatusError {
	return &StatusError{
		ErrStatus: Status{
			Message: message,
			Reason:  StatusReasonNotFound,
			Code:    http.StatusNotFound,
		},
	}
}

func NewUnknown(message string) *StatusError {
	return &StatusError{
		ErrStatus: Status{
			Message: message,
			Reason:  StatusReasonUnknown,
			Code:    -1,
		},
	}
}

func IsAPIError(err error) bool {
	_, ok := err.(APIStatus)
	return ok
}

func NewError(resp *http.Response) (*StatusError, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	apiError := &Error{}
	err = json.Unmarshal(body, apiError)
	if err != nil {
		return nil, err
	}

	var message string
	if apiError.Error != nil {
		message = *apiError.Error
	}

	if apiError.Text != nil {
		message = *apiError.Text
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return NewUnauthorized(message), nil
	case http.StatusForbidden:
		return NewForbidden(message), nil
	case http.StatusNotFound:
		return NewNotFound(message), nil
	default:
		return NewUnknown(message), nil
	}
}
