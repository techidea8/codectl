package logic

import (
	"fmt"
	"os"

	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/turingdance/infra/dbkit"
	"github.com/turingdance/infra/logger"
)

func InitApp(c *conf.AppConf) {
	level := logger.DebugLevel
	if c.Env == conf.PROD {
		level = logger.ErrorLevel
	}
	filewriter, err := dbkit.FileWriter(c.LogFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	engin, err := dbkit.OpenDb(dbkit.DBTYPE(c.DbType), c.Dsn,
		dbkit.WithWriter(os.Stdout, filewriter),
		dbkit.WithPrefix(c.Prefix),
		dbkit.IgnoreRecordNotFoundError(true),
		dbkit.ParameterizedQueries(true),
		dbkit.SetLogLevel(level),
		dbkit.SingularTable(true),
		dbkit.AutoMigrate(&model.Project{}, &model.Table{}, &model.Column{}),
	)
	if err != nil {
		panic(err)
	}
	DbEngin = engin
}
