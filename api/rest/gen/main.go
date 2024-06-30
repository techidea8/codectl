package main

import (
	"flag"

	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/infra/restkit"
	"github.com/techidea8/codectl/infra/restkit/middleware"
)

var env string

func main() {
	flag.StringVar(&env, "env", "dev", "envname: prod/dev")
	flag.Parse()
	c := conf.DefaultAppConf
	c.Env = conf.ENVDEF(env)
	logic.InitApp(c)
	restkit.ListenAndServe(":8080", restkit.NewHandler("/test").Pre(middleware.CORS(), middleware.NewAccessLog("access.log").Serve()))
}
