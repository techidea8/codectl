package biz

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/techidea8/codectl/infra/slicekit"
	"github.com/techidea8/codectl/infra/stringx"
)

var dirsrc string = ""
var dirdst string = ""

type Route struct {
	Module  string
	Func    string
	Path    string
	Method  []string
	Comment string
}

func (p *Route) Println() {
	fmt.Printf("moduel=%s,func=%s,path=%s,method=%s,comment=%s\n", p.Module, p.Func, p.Path, p.Method, p.Comment)
}

var rulemodule *regexp.Regexp = regexp.MustCompile(`\s*//\s+router\s+(\S+)`)
var rulepath *regexp.Regexp = regexp.MustCompile(`\s*//\s+((?:\,?(?:post|get|put|delete|options))+)\s+((?:/\w+)+)`)
var rulestruct *regexp.Regexp = regexp.MustCompile(`\s*type\s+(\S+)\s+struct\{.*\}`)

// var regfunc *regexp.Regexp = regexp.MustCompile(`.*func\s*\(\s*\w*\s*\*\s*(\w+)\s*\)\s*([\w]+)\s*\(\s*\S+\s*http\.ResponseWriter\s*\,\s*\S+\s*\*http\.Request\s*\).*`)
var rulefunc *regexp.Regexp = regexp.MustCompile(`\s*func\s+\(\s*\w+\s+\*([A-Z]+\w+)\s*\)\s+([A-Z]+\w+)\s*\(.*`)
var regcomment *regexp.Regexp = regexp.MustCompile(`.*\/\/\s*(.*).*`)

// 下划线单词转为小写驼峰单词
func camel(s string) string {

	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	ret := string(data[:])
	return strings.ToLower(ret[:1]) + ret[1:]
}
func buildroutes(dirsrc string) (routes []*Route, err error) {
	routes = make([]*Route, 0)
	err = filepath.WalkDir(dirsrc, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return err
		}
		bts, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(bytes.NewReader(bts))
		var proute *Route = nil
		var latestcomment string = ""
		for scanner.Scan() {
			txt := scanner.Text()
			if rulemodule.MatchString(txt) {
				// 模块
				// router /acc
				module := rulemodule.FindStringSubmatch(txt)
				proute = &Route{}
				proute.Path = module[1]

			} else if rulestruct.MatchString(txt) {
				// 结构体
				//
				structs := rulestruct.FindStringSubmatch(txt)
				if proute == nil {
					proute = &Route{}
					proute.Path = "/" + stringx.UnderlineToCamelCase(structs[1])
				}
				proute.Module = structs[1]
				proute.Comment = latestcomment
				routes = append(routes, proute)
				proute.Println()
				proute = nil
			} else if rulepath.MatchString(txt) {
				// 方法
				method := rulepath.FindStringSubmatch(txt)
				proute = &Route{}
				proute.Method = strings.Split(method[1], ",")
				proute.Path = method[2]
			} else if rulefunc.MatchString(txt) {
				result := rulefunc.FindStringSubmatch(txt)
				if proute == nil {
					proute = &Route{}
					proute.Method = []string{"post"}
					proute.Path = stringx.UnderlineToCamelCase("/" + stringx.CamelLcFirst(result[2]))
				}
				proute.Module = result[1]
				proute.Func = result[2]
				proute.Comment = latestcomment
				routes = append(routes, proute)
				proute.Println()
				proute = nil

			} else if regcomment.MatchString(txt) {
				latestcomment = regcomment.FindStringSubmatch(txt)[1]
			}
		}
		return err
	})

	return
}
func score(path string) int {
	var scorechar string = "{:"
	var ret int = 0
	for _, v := range scorechar {
		if strings.Contains(path, string(v)) {
			ret += 1
		}
	}
	return ret
}

var tplrouter string = `
//dot modify ! gen by  go run api/cmd/router/ -d api/rest/sys/handler -s  api/rest/sys/handler
package ${package}

import (
	"github.com/techidea8/codectl/infra/restkit"
)

var router *restkit.Router = restkit.NewRouter().PathPrefix("/")
// 初始化路由
func InitRouter(router *restkit.Router) {
	{{- range $k,$v := . }}
	{{$module := $v.Node.Module|camel}}
	// {{$v.Node.Comment}}
	{{$module}}Ctrl := &{{$v.Node.Module}}{}
	{{$module}}router := router.Subrouter().PathPrefix("{{$v.Node.Path}}")
	{{- range $g,$h := $v.Children }}
	//{{$h.Comment}}
	{{$module}}router.HandleFunc("{{$h.Path}}", {{$module}}Ctrl.{{$h.Func}}).Methods({{range $h.Method}}"{{.}}",{{end}})
	{{end}}

	{{end}}
}
func init() {
	InitRouter(router)
}
`

func replace(input string, rule map[string]string) string {
	for k, v := range rule {
		input = strings.ReplaceAll(input, k, v)
	}
	return input
}

type NodeRoute struct {
	Node     *Route
	Children []*Route
}

// 生成代码
func gencode(dirdst string, routes []*Route) (err error) {
	//NodeRoute
	noderoute := make([]*NodeRoute, 0)
	for _, v := range routes {
		if v.Func == "" {

			noderoute = append(noderoute, &NodeRoute{
				Node:     v,
				Children: make([]*Route, 0),
			})
		}
	}
	for _, v := range routes {
		if v.Func != "" {
			for _, v1 := range noderoute {
				if v1.Node.Module == v.Module {
					v1.Children = append(v1.Children, v)
				}
			}
		}
	}
	for _, v := range noderoute {
		slicekit.SortStable(v.Children, func(e1, e2 *Route) bool {
			score1 := score(e1.Path)
			score2 := score(e2.Path)
			return score1 > score2
		})
	}

	pkg := filepath.Base(dirdst)

	tpl, err := template.New("root").Funcs(template.FuncMap{"join": func(str []string) string {
		return strings.Join(str, ",")
	},
		"json": func(input any) string {
			bts, _ := json.Marshal(input)
			return string(bts)
		},
		"camel": camel,
	}).Parse(replace(tplrouter, map[string]string{
		"${package}": pkg,
	}))
	if err != nil {
		return err
	}
	dstfilename := path.Join(dirdst, "router.go")
	_, err = os.Stat(dstfilename)
	// 如果文件不存在,则创建
	if err == nil || os.IsNotExist(err) {
		err = os.Rename(dstfilename, dstfilename+".bak")
	}
	dstfile, err := os.Create(dstfilename)
	if err != nil {
		return err
	}
	dstfile.Truncate(0)
	defer dstfile.Close()
	err = tpl.Execute(dstfile, noderoute)
	return err

}

func gen(dirsrc string, dirdst string) error {
	routes, err := buildroutes(dirsrc)

	if err != nil {
		return err
	}
	return gencode(dirdst, routes)
}

// 子命令定义 运行方法 go run main.go version 编译后 ./hugo version
var routerCmd = &cobra.Command{
	Use:   "router", // Use这里定义的就是命令的名称
	Short: "通过注解生成路由",
	Long:  `generate route by annotation`,
	Run: func(cmd *cobra.Command, args []string) { //这里是命令的执行方法
		if dirsrc != "" && dirdst == "" {
			dirdst = dirsrc
		}
		//扫描目录下的每一个文件
		if err := gen(dirsrc, dirdst); err != nil {
			fmt.Println("gen route ❎", err.Error())
		} else {
			fmt.Println("gen route ✅")
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
	rootCmd.AddCommand(routerCmd)
	routerCmd.Flags().StringVarP(&dirsrc, "src", "s", "", "源目录")
	routerCmd.Flags().StringVarP(&dirdst, "dst", "d", "", "目标目录")
}
