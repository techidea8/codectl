package restkit

import (
	"github.com/techidea8/codectl/infra/wraper"
)

type (

	// HandlerFunc defines a function to serve HTTP requests.
	HandlerFunc func(ctx Context) (r *wraper.Response, err error)
	// MiddlewareFunc defines a function to process middleware.
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)
