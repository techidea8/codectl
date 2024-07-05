package handler

import (
	"fmt"
	"os"

	"github.com/techidea8/codectl/api/rest/gen/vo"
	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/techidea8/codectl/infra/cond"
	"github.com/techidea8/codectl/infra/dbkit"
	"github.com/techidea8/codectl/infra/filekit"
	"github.com/techidea8/codectl/infra/logger"
	"github.com/techidea8/codectl/infra/restkit"
	"github.com/techidea8/codectl/infra/wraper"
	"gorm.io/gorm"
)

type table struct{}

func init() {
	restkit.Register(&table{})
}

// 导出代码，生成zip文件
func (ctrl *table) Export(ctx restkit.Context) (r *wraper.Response, err error) {
	param := &vo.ExportVo{}
	err = ctx.Bind(param)
	if err != nil {
		return wraper.Error(err), err
	}
	//处理export
	prj, err := logic.TakeCurrentProject()
	if err != nil {
		return wraper.Error(err), err
	}
	if prj.Dsn == "" {
		return wraper.Error("dsn is empty"), fmt.Errorf("dsn is empty")
	}
	if prj.DbType == "" {
		return wraper.Error("data type is empty"), fmt.Errorf("data type is empty")
	}

	// 查询数据库
	// 解析prj.Dsn 获得databasename

	exportdb, err := dbkit.OpenDb(dbkit.DBTYPE(prj.DbType), prj.Dsn, dbkit.WithWriter(os.Stdout), dbkit.SetLogLevel(logger.ErrorLevel))

	if err != nil {
		logger.Error(err.Error(), "dsn", prj.Dsn)
		return
	}
	dbconf := dbkit.ParseMysql(prj.Dsn)
	dbname := dbconf.Dbname
	tablename := param.Name
	table := &model.Table{
		ProjectID: prj.ID,
		Name:      tablename,
	}
	// 取数据,有没有现成的
	table, err = logic.TakeTable(table)
	// 如果报错了 直接报错
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err.Error(), "logic.BuildFromSchema", prj.Dsn)
		return
	}
	tableexist := (table != nil && table.ID > 0)
	// 如果已经存在了,对数据进行添加和修改
	mapexistcolumn := make(map[string]model.Column, 0)
	if tableexist {
		// 获取clomun
		// 首先获得字段和row之间的map
		fields, _, err := logic.ListAllColumnByTableId(table.ID)
		if err != nil {
			return wraper.Error(err), err
		}
		for _, field := range fields {
			mapexistcolumn[field.RawData.ColumnName] = field
		}
	} else { // 如果表不存在,那么创建表结构
		table.Module = param.Module
		table.Name = param.Name
		table.Title = param.Title
		table.ProjectID = prj.ID
		table.Method = logic.BuildMethod(param.Methods)
		table, err = logic.CreateTable(table)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}
	// 新的数据列表,如果数据有变动则更新
	newcolumns, err := logic.BuildColumnFromSchema(exportdb, dbname, tablename)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	mapnewcolumn := make(map[string]model.Column, 0)
	keydbtypetolang := fmt.Sprintf("%s-%s", prj.DbType, prj.Lang)
	for _, newfield := range newcolumns {
		// 如果不存存在,说明作为添加存在
		if oldfield, exist := mapexistcolumn[newfield.RawData.ColumnName]; !exist {
			newfield.RawData.BuildColumn(&newfield)
			// 新的数据类型
			datatype := conf.DataTypeTOLangMapperRule[keydbtypetolang][newfield.RawData.DataType]
			if datatype == "" {
				datatype = newfield.RawData.DataType
			}
			newfield.DataType = datatype
			logger.Debug(newfield.RawData.DataType, "=>", datatype)
			mapnewcolumn[newfield.RawData.ColumnName] = newfield
		} else {
			oldfield.RawData = newfield.RawData
			oldfield.RawData.BuildColumn(&oldfield)
			// 获得实际数据类型映射
			datatype := conf.DataTypeTOLangMapperRule[keydbtypetolang][oldfield.RawData.DataType]
			if datatype == "" {
				datatype = newfield.RawData.DataType
			}
			logger.Debug(oldfield.RawData.DataType, "=>", datatype)
			oldfield.DataType = datatype
			mapexistcolumn[newfield.RawData.ColumnName] = oldfield
		}
	}
	tmpcreate := make([]model.Column, 0)
	tmpupdate := make([]model.Column, 0)
	// 继续保存到数据
	for _, field := range mapnewcolumn {
		field.TableID = table.ID
		tmpcreate = append(tmpcreate, field)
		//logic.DbEngin.Model(&model.Column{}).Create(&field)
	}
	if len(tmpcreate) > 0 {
		err := logic.DbEngin.Model(&model.Column{}).Create(tmpcreate).Error
		if err != nil {
			return wraper.Error(err), err
		}
	}
	//
	for _, field := range mapexistcolumn {
		logic.DbEngin.Model(&field).Updates(&field)
	}
	tmpcreate = append(tmpcreate, tmpupdate...)
	//kaishi export
	table.Columns = tmpcreate
	table.ProjectID = prj.ID
	table.Project = *prj
	table.Method = logic.BuildMethod(param.Methods)
	files := []string{}
	err = logic.ExportTable(table, conf.DirTpldata, func(file string) {
		files = append(files, file)
	})
	if err != nil {
		return wraper.Error(err), err
	}

	tmpfile, err := os.CreateTemp("/tmp", "codectltemp.*.zip")

	err = filekit.ZipFiles(tmpfile, files)
	if err != nil {
		return wraper.Error(err), err
	}
	filebytes := make([]byte, 0)
	_, err = tmpfile.Read(filebytes)
	if err != nil {
		return wraper.Error(err), err
	}
	defer tmpfile.Close()
	return wraper.Blob(wraper.BlobDef{
		File:        filebytes,
		Name:        param.Title + ".zip",
		ContentType: "application/zip",
	}).WithError(err), err

}

// 创建表格
func (prj *table) Create(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Table{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	instance, err = logic.CreateTable(instance)
	return wraper.OkData(instance).WithMsg("模块创建成功"), err
}

func (prj *table) Update(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Table{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	instance, err = logic.UpdateTable(instance, "id = ?", instance.ID)
	return wraper.OkData(instance).WithMsg("模块修改成功"), err
}
func (prj *table) List(ctx restkit.Context) (r *wraper.Response, err error) {
	prjs, total, err := logic.ListAllProject(&cond.CondWraper{
		Pager: cond.Pager{
			Pagesize: 1024,
		},
		Order: cond.Order{
			Field:  "sort_index",
			Method: "desc",
		},
	})
	return wraper.OkData(prjs).WithTotal(total), err
}

func (prj *table) Delete(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Table{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	total, err := logic.DeleteTable(instance, "id = ?", instance.ID)
	return wraper.OkData(total).WithMsg("模块删除成功"), err
}
