package webserver

import "net/http"

// ValidHttpStatusCodes is a list of valid status codes for the web server.
var ValidHttpStatusCodes = []int{
	http.StatusForbidden,
	http.StatusPreconditionFailed,
	http.StatusInternalServerError,
	http.StatusBadRequest,
	http.StatusBadGateway,
	http.StatusExpectationFailed,
	http.StatusTooEarly,
	http.StatusTooManyRequests,
	http.StatusServiceUnavailable,
	http.StatusMovedPermanently,
}
