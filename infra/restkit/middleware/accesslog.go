package middleware

import (
	"os"

	"github.com/techidea8/codectl/infra/logger"
	"github.com/techidea8/codectl/infra/restkit"
	"github.com/techidea8/codectl/infra/wraper"
)

type AccessLog struct {
	logger logger.ILogger
}

/*
$remote_addr             客户端地址                                    211.28.65.253
$remote_user             客户端用户名称                                --
$time_local              访问时间和时区                                18/Jul/2012:17:00:01 +0800
$request                 请求的URI和HTTP协议                           "GET /article-10000.html HTTP/1.1"
$http_host               请求地址，即浏览器中你输入的地址（IP或域名）     www.wang.com 192.168.100.100
$status                  HTTP请求状态                                  200
$upstream_status         upstream状态                                  200
$body_bytes_sent         发送给客户端文件内容大小                        1547
$http_referer            url跳转来源                                   https://www.baidu.com/
$http_user_agent         用户终端浏览器等信息                           "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; SV1; GTB7.0; .NET4.0C;
$ssl_protocol            SSL协议版本                                   TLSv1
$ssl_cipher              交换数据中的算法                               RC4-SHA
$upstream_addr           后台upstream的地址，即真正提供服务的主机地址     10.10.10.100:80
$request_time            整个请求的总时间                               0.205
$upstream_response_time  请求过程中，upstream响应时间                    0.002
————————————————
*/
func NewAccessLog(logpath string) *AccessLog {
	file, _ := os.OpenFile(logpath, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	logger := logger.NewStd(logger.InfoLevel, file, os.Stdout)
	return &AccessLog{
		logger: logger,
	}
}
func (a *AccessLog) SetLevel(level logger.LogLevel) *AccessLog {
	a.logger.SetLevel(level)
	return a
}
func (a *AccessLog) Serve() restkit.MiddlewareFunc {
	//
	return func(hf restkit.HandlerFunc) restkit.HandlerFunc {
		funcr := func(req restkit.Context) (r *wraper.Response, err error) {
			logger.Debugf("url=%s,ip=%s\n", req.Request().RequestURI, req.Request().RemoteAddr)
			return hf(req)
		}
		return funcr
	}
}
