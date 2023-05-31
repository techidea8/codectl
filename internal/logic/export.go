package logic

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"techidea8.com/codectl/internal/model"
	"techidea8.com/codectl/internal/util"
)

func prepare(col *model.Column) {
	if col.ColumnName == "" {
		col.ColumnName = util.UnderlineToCamelCase(col.DbColumnName)
	}
	if col.Title == "" {
		col.Title = col.ColumnName
	}
	col.IsCreateField = !col.IsPrimaryKey
	col.IsSearchField = !col.IsPrimaryKey
	col.IsPrimaryKey = strings.Compare(col.DbColumnKey, "PRI") == 0
	col.DataType = model.GuesDataType(col.DbDataType)
	col.DataMap = make(map[string]template.HTML)
	col.DataMap["Jsontag"] = col.BuildJsonTag()
	col.DataMap["Formtag"] = col.BuildFormTag()
	col.DataMap["Gormtag"] = col.BuildGormTag()
	col.Jsontag = col.BuildJsonTag()
	col.Formtag = col.BuildFormTag()
	col.Gormtag = col.BuildGormTag()
	fmt.Println(col.Jsontag, col.Formtag, col.Gormtag)
}
func Export(tab *model.Table) error {
	tab.Create()
	rows, _ := BuildCols(tab)
	mapRow := map[string]model.Column{}
	for _, v := range rows {
		prepare(&v)
		line := v.BuildLine()
		fmt.Println(line)
		mapRow[v.DbColumnName] = v
	}

	cols := model.Column{
		Database:  tab.Database,
		TableName: tab.TableName,
	}
	rows1, _ := cols.Search()
	ids := []uint{}
	for _, v := range rows1 {
		if ex, ok := mapRow[v.DbColumnName]; ok {
			ex.Id = v.Id
			mapRow[v.DbColumnName] = ex
		} else {
			ids = append(ids, v.Id)
		}
	}
	// 不存在的需要删除
	if len(ids) > 0 {
		model.NewColumn().Delete(ids)
	}
	for _, v := range mapRow {
		model.DbEngin.Save(&v)
	}
	rows1, _ = cols.Search()
	if err := Render(tab, rows1); err != nil {
		fmt.Println(err.Error())
	}
	return nil
	//return util.PrintTableWithAny(rows1)
}
func Render(tab *model.Table, columns []model.Column) (err error) {
	tab.Columns = make([]model.Column, 0)
	tab.Columns = append(tab.Columns, columns...)
	tmpls := template.New("root")
	tmpls = tmpls.Funcs(template.FuncMap{
		"ucfirst": util.Ucfirst,
		"lcfirst": util.Lcfirst,
		"jsstr":   util.JSStr,
		"js":      util.JS,
	})
	tmpls, err = tmpls.ParseGlob(tab.Tpldir + "/*")
	if err != nil {
		return err
	}
	for _, tpl := range tmpls.Templates() {
		tplName := tpl.Name()
		//过滤掉以html结尾的
		if strings.HasSuffix(tplName, ".html") {
			continue
		}
		//将
		dstFile := strings.ReplaceAll(tplName, "[model]", strings.ToLower(tab.Module))
		dstFile = strings.ReplaceAll(dstFile, "[Model]", tab.Module)
		pkgpath := strings.ReplaceAll(tab.Package, ".", "/")
		dstFile = strings.ReplaceAll(dstFile, "[pkgpath]", pkgpath)

		dstFile = filepath.Join(tab.Savedir, dstFile)
		dstFile = strings.TrimSuffix(dstFile, ".tpl")
		os.MkdirAll(filepath.Dir(dstFile), fs.FileMode(os.O_CREATE))

		f, e := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE, 0766)
		if e != nil {
			return e
		}

		//文件需要再次清空
		err = f.Truncate(0)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		err = tpl.ExecuteTemplate(f, tplName, *tab)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		f.Close()
		buf, _ := ioutil.ReadFile(dstFile)
		content := string(buf)
		content = strings.ReplaceAll(content, "&lt;", "<")
		err = ioutil.WriteFile(dstFile, []byte(content), 0766)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
	fmt.Println("generate code " + tab.TableName + "->" + tab.Module + " √")
	return os.Remove("root")
}
func listtables(input *model.Table) error {
	tmp := model.Table{
		Database: input.Database,
	}
	if rows, err := tmp.Search(); err != nil {
		fmt.Println(err.Error())
	} else {
		bts, _ := json.Marshal(rows)
		util.PrintTableWithByteArray(bts)
	}
	return nil
}
