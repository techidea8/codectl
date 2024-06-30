// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameProject = "project"

// Project mapped from table <project>
type Project struct {
	ID        int32  `gorm:"column:id;type:INTEGER;autoIncrement;primaryKey" json:"id"`
	Name     string `gorm:"column:name;type:string;size:50;" json:"name"`  //app name
	Title     string `gorm:"column:title;type:string;size:120;" json:"title"`
	DbType       string `gorm:"column:dbtype;type:string;size:30" json:"dbtype"`  //数据库类型
	DbName       string `gorm:"column:dbname;type:string;size:60" json:"dbname"`  //数据库名称
	Dsn       string `gorm:"column:dsn;type:string;size:250" json:"dsn"` //data source name
	Prefix       string `gorm:"column:prefix;type:string;size:50" json:"prefix"` //前缀
	Author       string `gorm:"column:author;type:string;size:50" json:"author"` //作者
	TplId    string `gorm:"column:tpl_id;type:string;size:120;" json:"tplId"`
	Package   string `gorm:"column:package;type:string;size:120;" json:"package"`
	SortIndex int32  `gorm:"column:sort_index;type:integer" json:"sortIndex"`
	Dirsave   string `gorm:"column:dirsave;type:string;size:120;" json:"dirsave"`
	Lang      string `gorm:"column:lang;type:string;size:20;default:golang" json:"lang"`
}

// TableName Project's table name
func (*Project) TableName() string {
	return TableNameProject
}
