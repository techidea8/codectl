package restkit

import (
	"net/http"

	"github.com/techidea8/codectl/infra/wraper"
)

type Transfer func(HandlerFunc) http.HandlerFunc

// 将我们的数据格式转化为标准的网络函数
func TransferHttpHandlerFunc(fun HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		resp, err := fun(ctx)
		if err != nil {
			wraper.Error(err.Error()).Encode(w)
		} else {
			resp.Encode(w)
		}
	}
}
