package logic

import (
	"techidea8.com/codectl/internal/model"
	"techidea8.com/codectl/internal/util"
)

// reverseCmd represents the reverse command
func Reverse(database string) error {
	ds, err := TakeDs(database)
	if err != nil {
		return err
	}
	tables, err := BuildTables(ds)
	if err != nil {
		return err
	}
	tabMap := map[string]model.Table{}
	for _, v := range tables {
		tabMap[v.TableName] = v
	}
	util.PrintTableWithAny(tables)
	//判断是否已经 存在到数据库里面
	table := &model.Table{Database: database}
	rows, err := table.Search()
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(tabMap))
	for k := range tabMap {
		keys = append(keys, k)
	}
	for _, v := range rows {
		// 如果不存在,就添加进去
		if exist, ok := tabMap[v.TableName]; ok {
			exist.Id = v.Id
			tabMap[v.TableName] = exist
		}
	}
	for _, v := range tabMap {
		model.DbEngin.Save(&v)
	}
	return nil
}
