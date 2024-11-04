package logic

import (
	"fmt"

	"github.com/techidea8/codectl/app/gen/model"
	"github.com/turingdance/infra/cond"
	"github.com/turingdance/infra/dbkit"
	"gorm.io/gorm"
)

// 列表全部项目
func ListAllTable(wraper *cond.CondWraper) (result []model.Table, total int64, err error) {
	wraper.Pager.Pagesize = 1024 * 1024
	return dbkit.Search(DbEngin, &model.Table{}, wraper)
}
func UpdateTable(model *model.Table, query interface{}, args ...interface{}) (*model.Table, error) {
	return dbkit.Update(DbEngin, model, query, args...)
}

// 获得当前项目
func TakeTableById(tableId int32) (result *model.Table, err error) {
	result = &model.Table{}
	err = DbEngin.Where("ID = ?", tableId).Preload("Project").Preload("Columns", " table_id = ?", tableId).Find(result).Error
	return result, err
}

// 获得当前项目
func TakeTable(instance *model.Table) (result *model.Table, err error) {
	result = &model.Table{}
	err = DbEngin.Model(instance).Where(instance).Preload("Project").Preload("Columns").Find(result).Error
	return result, err
}

// 获得当前项目
func MetaFromModuleName(module string) (result []model.Column, err error) {
	table := &model.Table{}
	col := &model.Column{}
	result = make([]model.Column, 0)
	err = DbEngin.Model(result).Where("name = ?", module).Find(table).Error
	if table != nil && err == nil {
		err = DbEngin.Model(col).Where("table_id = ?", table.ID).Find(result).Error
	}
	return result, err
}

func CreateTable(model *model.Table) (*model.Table, error) {
	return dbkit.Create(DbEngin, model)
}

func DeleteTable(model *model.Table, query interface{}, args ...interface{}) (total int64, err error) {
	return dbkit.Delete(DbEngin, model, query, args...)
}

// 从数据库中获取table信息
func BuildTableFromSchema(dbengin *gorm.DB, dbname string) (comumns []model.Table, err error) {
	dbtype := dbengin.Dialector.Name()
	if dbtype == "mysql" {
		return BuildTableFromMysqlSchema(dbengin, dbname)
	} else {
		return []model.Table{}, fmt.Errorf("不支持的数据库类型" + dbtype)
	}
}
func BuildTableFromMysqlSchema(dbengin *gorm.DB, dbname string) (comumns []model.Table, err error) {
	tables := make([]model.Table, 0)
	err = dbengin.Raw(`select TABLE_NAME as name,TABLE_COMMENT as title
	from information_schema.tables where  table_schema = ? `, dbname).Scan(&tables).Error

	return tables, err
}
