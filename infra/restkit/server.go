package restkit

import (
	"net/http"
	"time"
)

// 监听服务
func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)

}

type ServerOpt func(*http.Server)

func ReadTimeout(tm time.Duration) ServerOpt {
	return func(s *http.Server) {
		s.ReadTimeout = tm
	}
}
func WriteTimeout(tm time.Duration) ServerOpt {
	return func(s *http.Server) {
		s.WriteTimeout = tm
	}
}

func IdleTimeout(tm time.Duration) ServerOpt {
	return func(s *http.Server) {
		s.IdleTimeout = tm
	}
}

func NewServer(addr string, hander http.Handler, opts ...ServerOpt) *http.Server {
	svr := &http.Server{
		Addr:         addr,
		Handler:      hander,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  0,
	}
	for _, v := range opts {
		v(svr)
	}
	return svr
}
