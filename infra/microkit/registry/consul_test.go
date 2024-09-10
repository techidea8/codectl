package registry_test

import (
	"net/http"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"github.com/techidea8/codectl/infra/microkit/registry"
)

func Test01(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	c := consulapi.DefaultConfig()
	r := registry.NewConsulRegistry(c)
	svc := &registry.Service{
		Name: "test",
		Host: "127.0.0.1",
		Port: 1234,
		Tags: []string{
			"traefik.http.routers.turing-account.rule=PathPrefix(`/account`)",
			"traefik.http.routers.turing-account.middlewares=turing-replacepath,turing-staticserver,turing-jwtauthor,turing-compress",
			"traefik.http.middlewares.turing-replacepath.stripPrefix.prefixes=/account",
			"traefik.http.middlewares.turing-staticserver.plugin.traefikstaticfs.alias=/mnt:d:/data/storage/mnt,/asset:e:/data/asset",
			"traefik.http.middlewares.turing-jwtauthor.plugin.traefikjwtauthor.whiteList=/login,/mnt,/register",
			"traefik.http.middlewares.turing-jwtauthor.plugin.traefikjwtauthor.secret=secret",
			"traefik.http.middlewares.turing-compress.compress=true",
		},
		Check: registry.Check{
			DeregisterCriticalServiceAfter: 1 * time.Second,
		},
	}
	r.Register(svc)
	e.Logger.Fatal(e.Start(":1234"))
}
