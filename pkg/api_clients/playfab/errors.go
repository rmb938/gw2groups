package playfab

import "fmt"

type APIErrorWrapper interface {
	GetErrorCode() int
}

type APIErrorResponse struct {
	APIResponse  `json:"-"`
	AError       string                 `json:"error"`
	ErrorCode    int                    `json:"errorCode"`
	ErrorDetails map[string]interface{} `json:"errorDetails"`
	ErrorMessage string                 `json:"errorMessage"`
}

func (e *APIErrorResponse) GetErrorCode() int {
	return e.ErrorCode
}

func (e *APIErrorResponse) Error() string {
	return fmt.Sprintf("received error from playfab: %d(%d - %s) %s: %#v", e.Code, e.ErrorCode, e.AError, e.ErrorMessage, e.ErrorDetails)
}

func IsAPIError(err error) bool {
	_, ok := err.(APIErrorWrapper)
	return ok
}
