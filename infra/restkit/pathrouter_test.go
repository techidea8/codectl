package restkit

import (
	"fmt"
	"testing"

	"github.com/techidea8/codectl/infra/wraper"
)

func Test01(t *testing.T) {
	router := NewRouter().PathPrefix("hello")
	testrouter := router.Subrouter().PathPrefix("/test")
	testrouter.HandleFunc("test01", func(ctx Context) (r *wraper.Response, err error) {
		return wraper.Error("1"), nil
	}).Methods("POST", "GET")
	testrouter.HandleFunc("/test02", func(ctx Context) (r *wraper.Response, err error) {
		return wraper.Error("1"), nil
	}).Methods("POST", "GET")
	testrouter.HandleFunc("test03", func(ctx Context) (r *wraper.Response, err error) {
		return wraper.Error("1"), nil
	}).Methods("POST", "GET")
	fmt.Println(router)
}
