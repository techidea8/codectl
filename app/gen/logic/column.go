package logic

import (
	"fmt"

	"github.com/techidea8/codectl/app/gen/model"
	"github.com/techidea8/codectl/infra/cond"
	"github.com/techidea8/codectl/infra/dbkit"
	"gorm.io/gorm"
)

// 列表全部项目
func ListAllColumn(wraper *cond.CondWraper) (result []model.Column, total int64, err error) {
	wraper.Pager.Pagesize = 1024 * 1024
	return dbkit.Search(DbEngin, &model.Column{}, wraper)
}
func ListAllColumnByTableId(tableId int32) (result []model.Column, total int64, err error) {
	return dbkit.Search(DbEngin, &model.Column{}, cond.NewCondWrapper().Pagesize(10086).AddCond(cond.Cond{
		Field: "table_id",
		Op:    cond.OPEQ,
		Value: tableId,
	}))
}
func CreateColumn(model *model.Column) (*model.Column, error) {
	return dbkit.Create(DbEngin, model)
}
func UpdateColumn(model *model.Column, query interface{}, args ...interface{}) (*model.Column, error) {
	return dbkit.Update(DbEngin, model, query, args...)
}
func DeleteColumn(model *model.Column, query interface{}, args ...interface{}) (total int64, err error) {
	return dbkit.Delete(DbEngin, model, query, args...)
}

func BuildColumnFromSchema(dbengin *gorm.DB, dbname, tablename string) (comumns []model.Column, err error) {
	dbtype := dbengin.Dialector.Name()
	if dbtype == "mysql" {
		return BuildColumnFromMysqlSchema(dbengin, dbname, tablename)
	} else {
		return []model.Column{}, fmt.Errorf("不支持的数据库类型" + dbtype)
	}
}
func BuildColumnFromMysqlSchema(dbengin *gorm.DB, dbname, tablename string) (comumns []model.Column, err error) {
	rawdatas := make([]model.RawData, 0)
	err = dbengin.Raw(`select COLUMN_NAME as column_name,
	DATA_TYPE as data_type,
	IFNULL(CHARACTER_MAXIMUM_LENGTH,0) as char_max_length,
	IFNULL(NUMERIC_PRECISION,0) as number_precision,
	IFNULL(NUMERIC_SCALE,0)  as number_scale,
	COLUMN_COMMENT as column_comment,
	IS_NULLABLE as is_nullable,
	COLUMN_KEY as column_key,
	COLUMN_TYPE as column_type,
	EXTRA as extra,
	ORDINAL_POSITION as ordinal_position  
	from information_schema.columns where  table_schema = ? and  table_name = ? order by ORDINAL_POSITION asc`, dbname, tablename).Scan(&rawdatas).Error
	comumns = make([]model.Column, 0)
	for _, raw := range rawdatas {
		comumn := &model.Column{}
		comumn = raw.BuildColumn(comumn)
		comumns = append(comumns, *comumn)
	}
	return comumns, err
}
