package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/viper"
	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/infra/restkit"
)

var (
	addr string
)

type config struct {
	gen *conf.AppConf
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	flag.StringVar(&addr, "addr", ":2048", "监听端口")
	flag.Parse()
	vp := viper.New()
	vp.SetConfigFile("app.yml")
	cfg := &config{
		gen: &conf.AppConf{},
	}
	vp.Unmarshal(cfg)
	logic.InitApp(cfg.gen)
	go restkit.ListenAndServe(addr, restkit.NewHandler("/"))
	fmt.Println("print ctrl+c to quit")
	for {
		select {
		case <-c:
			fmt.Println("quit app")
			return
		}
	}
}
