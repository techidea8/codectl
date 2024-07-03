package biz

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra" // 安装依赖 go get -u github.com/spf13/cobra/cobra
	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/techidea8/codectl/infra/cond"
)

type prjctrl struct {
}

// 列表
func (s *prjctrl) list(args []string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "服务名", "开发者", "应用描述", "数据库类型", "数据库名称", "数据库连接串", "表结构前缀", "模板", "应用包名", "存储位置", "语言类型", "排序位"})
	prjs, total, err := logic.ListAllProject(&cond.CondWraper{
		Pager: cond.Pager{
			Pagesize: 1024,
		},
		Order: cond.Order{
			Field:  "sort_index",
			Method: "desc",
		},
	})
	if err != nil {
		return err
	}

	for _, prj := range prjs {
		table.Append([]string{
			strconv.Itoa(int(prj.ID)),
			prj.Name,
			prj.Author,
			prj.Title,
			prj.DbType,
			prj.DbName,
			prj.Dsn,
			prj.Prefix,
			prj.TplId,
			prj.Package,
			prj.Dirsave,
			prj.Lang,
			fmt.Sprintf("%d", prj.SortIndex)})
	}
	fmt.Printf("num  = %d\n", total)
	table.Render()
	return nil
}

// 添加
func (s *prjctrl) add(args []string) error {
	prj := &model.Project{}
	s.bindprj("please input data ", prj)
	if prj.Package == "" {
		prj.Package = conf.DefaultPackage
	}
	if prj.TplId == "" {
		prj.Package = conf.DefaultTplId
	}
	if prj.Name == "" {
		prj.Name = "biz"
	}
	if prj.Lang == "" {
		prj.Lang = conf.DefaultLang
	}
	prj.SortIndex = int32(time.Now().Unix())
	logic.Create(prj)
	s.list(args)
	return nil
}
func (s *prjctrl) bindprj(title string, prj *model.Project) error {
	fmt.Println(title)
	fmt.Println("name,the name of project,must be english ,will use as service name ")
	fmt.Print(">")
	fmt.Scanln(&prj.Name)

	fmt.Println("author , author of this app ")
	fmt.Print(">")
	fmt.Scanln(&prj.Author)

	fmt.Println("title,the title of project")
	fmt.Print(">")
	fmt.Scanln(&prj.Title)

	//
	fmt.Println("dbtype,required eg:mysql/sqlite/postgres/sqlserver")
	fmt.Print(">")
	fmt.Scanln(&prj.DbType)

	//
	fmt.Println("dbname,required eg: database name ")
	fmt.Print(">")
	fmt.Scanln(&prj.DbName)

	fmt.Println("dsn,required eg:appdata.db or user:password@host:port/db?quest>")
	fmt.Print(">")
	fmt.Scanln(&prj.Dsn)

	fmt.Println("prefix,prefix of table,such as sys_, or kf_>")
	fmt.Print(">")
	fmt.Scanln(&prj.Prefix)

	fmt.Println("package,eg is turingdance.com/turing/app>")
	fmt.Print(">")
	fmt.Scanln(&prj.Package)

	fmt.Println("tplId,the data from tpl list")
	fmt.Print(">")
	fmt.Scanln(&prj.TplId)

	fmt.Println("dir for save,the path for saving code  such as path/to/tpl")
	fmt.Print(">")
	fmt.Scanln(&prj.Dirsave)

	fmt.Println("language eg:golang/java,defaultis golang")
	fmt.Print(">")
	fmt.Scanln(&prj.Lang)
	return nil
}
func (s *prjctrl) update(args []string) error {
	prjId, _ := strconv.Atoi(args[0])
	prj := &model.Project{
		ID: int32(prjId),
	}
	s.bindprj("please input data ,empty will not modify", prj)
	prj.SortIndex = int32(time.Now().Unix())
	logic.Update(prj, "id = ?", prj.ID)
	s.list(args)
	return nil
}

func (s *prjctrl) del(args []string) (err error) {
	prjId, _ := strconv.Atoi(args[0])
	prj := &model.Project{
		ID: int32(prjId),
	}
	_, err = logic.Delete(prj, "id = ?", prj.ID)
	s.list(args)
	return err
}
func (s *prjctrl) active(args []string) (err error) {
	prjId, _ := strconv.Atoi(args[0])
	prj := &model.Project{
		ID:        int32(prjId),
		SortIndex: int32(time.Now().Unix()),
	}
	_, err = logic.Update(prj, "id = ?", prjId)
	s.list(args)
	return err
}

// prjId
func (s *prjctrl) clone(args []string) (err error) {
	prjId, _ := strconv.Atoi(args[0])
	prj := &model.Project{
		ID: int32(prjId),
	}
	prj, err = logic.TakeByPrimaryKey(prj)
	prj.ID = 0
	logic.CreateProject(prj)
	s.list(args)
	return err
}

var defaultprjctrl *prjctrl = &prjctrl{}
var projectmapfun map[string]func([]string) error = map[string]func([]string) error{
	"list":   defaultprjctrl.list,
	"add":    defaultprjctrl.add,
	"update": defaultprjctrl.update,
	"del":    defaultprjctrl.del,
	"clone":  defaultprjctrl.clone,
	"active": defaultprjctrl.active,
	"use":    defaultprjctrl.active,
}

// 子命令定义 运行方法 go run main.go version 编译后 ./hugo version
var projectCmd = &cobra.Command{
	Use:   "project", // Use这里定义的就是命令的名称
	Short: "project manager",
	Long: `
project list
    list all project
project add
    add a project step by step
project update [projectId]
    update a project step by step,eg:project update 1
project del [projectId]
    delete a project ,eg: project del 1
project clone [projectId]
    clone a project ,eg: project clone 1
project use [projectId]
    set project as default ,eg project use 1
`,
	Run: func(cmd *cobra.Command, args []string) { //这里是命令的执行方法
		if len(args) < 1 {

		} else {
			args = append(args, "")
			projectmapfun[args[0]](args[1:])
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行前执行
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行后执行
	},
	// 还有其他钩子函数
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
