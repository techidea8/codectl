package logic

import (
	"techidea8.com/codectl/internal/model"
	"techidea8.com/codectl/internal/util"
)

func BuildTables(ds model.Ds) (tables []model.Table, err error) {
	db, err := ds.BuildDbEngin()
	if err != nil {
		return
	}
	tables = make([]model.Table, 0)
	err = db.Raw(`select table_name as table_name,table_comment as title from information_schema.tables where table_schema=?`, ds.Database).Scan(&tables).Error
	for i, v := range tables {
		tables[i].Database = ds.Database
		tables[i].Module = util.UnderlineToCamelCase(v.TableName)
		tables[i].Method = append([]model.Method{}, model.DefaultMethods...)
	}
	return tables, err
}
