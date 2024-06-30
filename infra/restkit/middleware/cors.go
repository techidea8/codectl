package middleware

import (
	"github.com/techidea8/codectl/infra/restkit"
	"github.com/techidea8/codectl/infra/wraper"
)

// CORS returns a Cross-Origin Resource Sharing (CORS) middleware.
// See: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
func CORS() restkit.MiddlewareFunc {
	return func(next restkit.HandlerFunc) restkit.HandlerFunc {
		return func(c restkit.Context) (r *wraper.Response, err error) {
			// 允许特定的域进行跨域请求
			c.Writer().Header().Set("Access-Control-Allow-Origin", "*")
			// 允许特定的请求方法
			c.Writer().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			// 允许特定的请求头
			c.Writer().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			// 允许携带身份凭证（如Cookie）
			c.Writer().Header().Set("Access-Control-Allow-Credentials", "true")
			// 继续处理请求
			return next(c)
		}
	}
}
