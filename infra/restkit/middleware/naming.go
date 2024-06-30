package middleware

import (
	"fmt"

	"github.com/techidea8/codectl/infra/restkit"
	"github.com/techidea8/codectl/infra/wraper"
)

// CORS returns a Cross-Origin Resource Sharing (CORS) middleware.
// See: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
func Naming(name string) restkit.MiddlewareFunc {
	return func(next restkit.HandlerFunc) restkit.HandlerFunc {
		return func(c restkit.Context) (r *wraper.Response, err error) {
			fmt.Printf("middleware %s\n", name)
			return next(c)
		}
	}
}
