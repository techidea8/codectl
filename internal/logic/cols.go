package logic

import "techidea8.com/codectl/internal/model"

func BuildCols(tab *model.Table) (cols []model.Column, err error) {
	ds := tab.CreateDs()
	db, err := ds.BuildDbEngin()
	if err != nil {
		return
	}

	rows := make([]model.Column, 0)
	err = db.Raw(`select COLUMN_NAME as db_column_name,
	DATA_TYPE as db_data_type,
	IFNULL(CHARACTER_MAXIMUM_LENGTH,0) as char_max_len,
	COLUMN_TYPE,IFNULL(NUMERIC_PRECISION,0) as nump,
	IFNULL(NUMERIC_SCALE,0)  as nump,
	COLUMN_COMMENT as title,
	column_key as db_column_key,
	column_type as db_column_type,
	extra as extra,
	ORDINAL_POSITION as ordinal_position  
	from information_schema.COLUMNS where  table_schema = ? and  table_name = ?`, tab.Database, tab.TableName).Scan(&rows).Error
	for i, _ := range rows {
		rows[i].Database = ds.Database
		rows[i].TableName = tab.TableName
	}
	return rows, err
}
