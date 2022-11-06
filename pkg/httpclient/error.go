package httpclient

import "fmt"

type ErrHTTPClient struct {
	StatusCode int
	msg        string
}

func (e *ErrHTTPClient) Error() string {
	return fmt.Sprintf("HTTPClient_Error[%d]: %s", e.StatusCode, e.msg)
}

func NewErrHTTPClient(statusCode int, msg string) error {
	return &ErrHTTPClient{
		StatusCode: statusCode,
		msg:        msg,
	}
}
