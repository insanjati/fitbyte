package errors

import "errors"

// 4xx Client Errors
var (
	ErrBadRequest          = errors.New("bad request")                   // 400
	ErrUnauthorized        = errors.New("unauthorized")                  // 401
	ErrPaymentRequired     = errors.New("payment required")              // 402
	ErrForbidden           = errors.New("forbidden")                     // 403
	ErrNotFound            = errors.New("not found")                     // 404
	ErrMethodNotAllowed    = errors.New("method not allowed")            // 405
	ErrNotAcceptable       = errors.New("not acceptable")                // 406
	ErrProxyAuthRequired   = errors.New("proxy authentication required") // 407
	ErrRequestTimeout      = errors.New("request timeout")               // 408
	ErrConflict            = errors.New("conflict")                      // 409
	ErrGone                = errors.New("gone")                          // 410
	ErrLengthRequired      = errors.New("length required")               // 411
	ErrPreconditionFailed  = errors.New("precondition failed")           // 412
	ErrPayloadTooLarge     = errors.New("payload too large")             // 413
	ErrURITooLong          = errors.New("uri too long")                  // 414
	ErrUnsupportedMedia    = errors.New("unsupported media type")        // 415
	ErrRangeNotSatisfiable = errors.New("range not satisfiable")         // 416
	ErrExpectationFailed   = errors.New("expectation failed")            // 417
	ErrTooManyRequests     = errors.New("too many requests")             // 429
)

// 5xx Server Errors
var (
	ErrInternalServerError   = errors.New("internal server error")      // 500
	ErrNotImplemented        = errors.New("not implemented")            // 501
	ErrBadGateway            = errors.New("bad gateway")                // 502
	ErrServiceUnavailable    = errors.New("service unavailable")        // 503
	ErrGatewayTimeout        = errors.New("gateway timeout")            // 504
	ErrHTTPVersionNotSupport = errors.New("http version not supported") // 505
)
