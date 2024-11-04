package biz

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra" // 安装依赖 go get -u github.com/spf13/cobra/cobra
	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/turingdance/infra/filekit"
	"github.com/turingdance/infra/logger"
	"github.com/turingdance/infra/oskit"
)

type tplctrl struct {
	tpldir string
}

func (s *tplctrl) walkdir(filepath string) ([]string, error) {
	files, err := os.ReadDir(filepath) // files为当前目录下的所有文件名称【包括文件夹】
	if err != nil {
		return []string{}, err
	}
	var allfile []string
	for _, v := range files {
		allfile = append(allfile, v.Name())
	}

	return allfile, nil
}

// 列表
func (s *tplctrl) list(args []string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"序号", "模板ID"})
	dirs, err := s.walkdir(s.tpldir)
	total := len(dirs)
	for index, tplId := range dirs {
		table.Append([]string{
			strconv.Itoa(index + 1),
			tplId,
		})
	}
	fmt.Printf("tpldir = %s\ntotal = %d\n", s.tpldir, total)
	table.Render()
	return err
}

// 添加
func (s *tplctrl) add(args []string) (err error) {
	src := args[0]
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return s.add_net(args)
	} else {
		return s.add_local(args)
	}
}

func (s *tplctrl) add_local(args []string) (err error) {
	src := args[0]

	total, err := filekit.Copy(src, s.tpldir)
	fmt.Printf("copy %s to %s , num = %d", src, s.tpldir, total)
	if err != nil {
		return err
	}
	return s.list(args)
}

// 添加网络模版
func (s *tplctrl) add_net(args []string) (err error) {
	dstdir := ""
	if len(args) > 1 {
		dstdir = args[1]
	} else {
		urls, err := url.Parse(args[0])
		if err != nil {
			return err
		}
		paths := strings.Split(urls.Path, "/")
		dstname := paths[len(paths)-1]
		dstdir = strings.TrimSuffix(dstname, ".git")
	}
	dstdir = strings.ReplaceAll(dstdir, "/", "_")
	dstdir = strings.ReplaceAll(dstdir, "\\", "_")
	dstdir = path.Join(s.tpldir, dstdir)
	return s.clone([]string{args[0], dstdir})
}

// 添加
func (s *tplctrl) use(args []string) (err error) {
	prj := &model.Project{}
	tplId := args[0]
	if prjId == 0 {
		prj, err = logic.TakeCurrentProject()
	} else {
		prj.ID = prjId
	}
	//
	prj.TplId = tplId
	logic.Update(prj, "ID = ?", prj.ID)
	(&prjctrl{}).list(args)
	return err
}

// 添加
func (s *tplctrl) del(args []string) error {
	prj := &model.Project{}

	prj.SortIndex = int32(time.Now().Unix())
	logic.Create(prj)
	s.list(args)
	return nil
}

// 添加
func (s *tplctrl) clone(args []string) error {
	argcmd := []string{
		"clone",
	}
	argcmd = append(argcmd, args...)
	resultch, errch, stopch := s.git(argcmd)
	for {
		select {
		case m := <-resultch:
			logger.Debug(m)
		case e := <-errch:
			logger.Debug(e)
		case e := <-stopch:
			return errors.New(e)
		}
	}
}
func (s *tplctrl) git(args []string) (resultch, errorch, stopch chan string) {
	ctx := context.Background()
	return oskit.ExecWithChanel(ctx, "git", args)
}
func NewtplCtl() *tplctrl {
	return &tplctrl{
		tpldir: conf.DirTpldata,
	}
}

var defaulttplctrl *tplctrl = NewtplCtl()
var tplmapfun map[string]func([]string) error = map[string]func([]string) error{
	"list":  defaulttplctrl.list,
	"add":   defaulttplctrl.add,
	"use":   defaulttplctrl.use,
	"del":   defaulttplctrl.del,
	"clone": defaulttplctrl.clone,
}

// 子命令定义 运行方法 go run main.go version 编译后 ./hugo version
var tplCmd = &cobra.Command{
	Use:   "tpl", // Use这里定义的就是命令的名称
	Short: "tpl manager",
	Long: `
tpl list
    list all template exist in tpl dir  
tpl use [tplid]
    use tplid as the default tpl for current project
tpl del [tplid]
    delete tpl
tpl add [tplpath] [dstdir]
    if tplpath is a net url,download it and make it as a templete
	if tplpath is a file url,copy it and make it as a templete
tpl clone [giturl] [dstdir]
    clone giturl  and save it at current dir for eg  tpl  clone https://github.com/techidea8/tpl-vue3-go.git
`,
	Run: func(cmd *cobra.Command, args []string) { //这里是命令的执行方法
		if len(args) == 0 {
			cmd.Usage()
			return
		}
		if len(args) == 1 {
			args = append(args, "")
		}
		//args = append(args, "", "")
		cmdstr := args[0]

		paramargs := args[1:]
		if fun, ok := tplmapfun[cmdstr]; ok {
			if err := fun(paramargs); err != nil {
				logger.Error(err.Error())
			}
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
var (
	prjId int32 = 0
)

func init() {
	rootCmd.AddCommand(tplCmd)
	tplCmd.Flags().Int32VarP(&prjId, "prj", "p", 0, "id of project ")
}
