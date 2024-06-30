package conf

import (
	"os"
	"path"
)

type ENVDEF string

const PROD ENVDEF = "prod"
const DEV ENVDEF = "dev"

type AppConf struct {
	DbType  string
	Dsn     string
	Env     ENVDEF
	LogFile string
	Prefix  string
}

const (
	DefaultTplId   = "turing-vue3-go-v2"
	DefaultPackage = "turngdance.com/turing/app"
	DefaultLang    = "golang"
)
const _dirappconfig = ".codectl"

var DirUserConfig, _ = os.UserConfigDir()
var (
	DirAppConfig = path.Join(DirUserConfig, _dirappconfig)
	DirTpldata   = path.Join(DirUserConfig, _dirappconfig, "tpldir")
)
var DefaultAppConf *AppConf = &AppConf{
	DbType:  "sqlite",
	Dsn:     path.Join(DirAppConfig, "codectl.db"),
	Env:     PROD,
	Prefix:  "gen_",
	LogFile: "app.log",
}

func init() {
	os.MkdirAll(DirTpldata, os.ModeDir)
}
