package apierror

import (
	"fmt"
	"net/http"
)

// CheckResponse checks the response for errors and returns the value or an error
func CheckResponse[T any](resp APIResponse[T], err error) (*T, error) {
	if err != nil {
		return nil, err
	}

	value := resp.GetValue()
	if value != nil {
		return value, nil
	}

	//nolint:bodyclose
	return nil, WithRequestID(responseError(resp), resp.GetHTTPResponse().Header)
}

type APIResponse[T any] interface {
	GetHTTPResponse() *http.Response
	GetInternalServerError() (string, bool)
	GetBadRequestError() (string, bool)
	GetValue() *T
}

type Error interface {
	GetError() string
}

func responseError[T any](resp APIResponse[T]) error {
	reason, ok := resp.GetBadRequestError()
	if ok {
		return NewTerminalError(fmt.Errorf("bad request: %v", reason))
	}
	reason, ok = resp.GetInternalServerError()
	if ok {
		return fmt.Errorf("internal server error: %v", reason)
	}
	//nolint:bodyclose
	httpResp := resp.GetHTTPResponse()
	if httpResp == nil {
		return fmt.Errorf("unexpected response: nil")
	}
	return fmt.Errorf("unexpected response: %v", httpResp.StatusCode)
}
