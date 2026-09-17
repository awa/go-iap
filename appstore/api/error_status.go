package api

import "fmt"

// httpStatusError carries the status of a response that arrived with no App
// Store error code. Its message is the one StoreClient always used. The type
// stays unexported; callers reach the status through StatusCode:
//
//	var withStatus interface{ StatusCode() int }
//	if errors.As(err, &withStatus) && withStatus.StatusCode() == http.StatusTooManyRequests {
type httpStatusError struct {
	statusCode int
	url        string
}

func newHTTPStatusError(statusCode int, url string) *httpStatusError {
	return &httpStatusError{statusCode: statusCode, url: url}
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("appstore api: %v return status code %v", e.url, e.statusCode)
}

// StatusCode is the status the host answered with.
func (e *httpStatusError) StatusCode() int {
	return e.statusCode
}
