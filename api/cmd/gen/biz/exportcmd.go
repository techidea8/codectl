package biz

import (
	"fmt"
	"os"

	"github.com/spf13/cobra" // 安装依赖 go get -u github.com/spf13/cobra/cobra
	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/turingdance/infra/dbkit"
	"github.com/turingdance/infra/logger"
)

type exportctrl struct {
	table   *model.Table // = &model.Table{}
	methods []string     //= make([]string, 0)
}

func NewExportCtrl() *exportctrl {
	return &exportctrl{
		table:   &model.Table{},
		methods: make([]string, 0),
	}
}
func (s *exportctrl) Init() {
	s.methods = make([]string, 0)
	s.table = &model.Table{}
}

// 子命令定义 运行方法 go run main.go version 编译后 ./hugo version
var exportCmd = &cobra.Command{
	Use:   "export", // Use这里定义的就是命令的名称
	Short: "export table to code",
	Long:  `export table to code such as golang/java..`,
	Run: func(cmd *cobra.Command, args []string) { //这里是命令的执行方法
		// 重置规则
		if rulefile != "" {
			conf.ResetMapperRule(rulefile)
		}
		//处理export
		prj, err := logic.TakeCurrentProject()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if prj.Dsn == "" {
			fmt.Print("please check dsn")
			return
		}
		if prj.DbName == "" {
			fmt.Print("please check data base config use project list")
			return
		}
		if prj.DbType == "" {
			fmt.Print("please check data base config use project list")
			return
		}

		// 查询数据库
		// 解析prj.Dsn 获得databasename
		loglevel := logger.InfoLevel
		if env == string(conf.PROD) {
			loglevel = logger.ErrorLevel
		}
		exportdb, err := dbkit.OpenDb(dbkit.DBTYPE(prj.DbType), prj.Dsn, dbkit.WithWriter(os.Stdout),
			dbkit.SetLogLevel(loglevel))
		if err != nil {
			logger.Error(err.Error())
			return
		}
		vo := &logic.PrepareVo{
			Project:    prj,
			BizDbEngin: exportdb,
			TableName:  defaultexportctrl.table.Name,
			ModuleName: defaultexportctrl.table.Module,
			BizTitle:   defaultexportctrl.table.Title,
			Methods:    model.SimpleMethods,
		}
		table, err := logic.PrepareExportTable(vo)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		err = logic.ExportTable(table, conf.DirTpldata)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行前执行
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行后执行
		defaultexportctrl.Init()
	},
	// 还有其他钩子函数
}
var defaultexportctrl *exportctrl = NewExportCtrl()

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&(defaultexportctrl.table.Name), "name", "t", "", "module name of table")
	exportCmd.Flags().StringVarP(&(defaultexportctrl.table.Title), "title", "n", "", "title of table")
	exportCmd.Flags().StringVarP(&(defaultexportctrl.table.Module), "module", "m", "", "module of table")
	exportCmd.Flags().StringVarP(&(rulefile), "rule", "r", "", "rule config file,like mapper.yml")
	exportCmd.Flags().StringArrayVarP(&defaultexportctrl.methods, "method", "f", []string{}, "method of table:search/create/update/take/download/export eg -f search -f create")
}
