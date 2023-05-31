package model

import (
	"html/template"
	"strings"

	"techidea8.com/codectl/internal/util"
)

type Column struct {
	Id              uint `json:"id" form:"id" gorm:"primaryKey;int(11) auto_increment"`
	Database        string
	TableName       string
	DbColumnName    string `json:"dbColumnName" form:"dbColumnName"`
	DbDataType      string `json:"dbDataType" form:"dbDataType"`
	DbColumnKey     string 
	ColumnName      string `json:"columnName" form:"columnName"`
	DbColumnType    string
	DataType        string //string
	CharMaxLen      int    //20
	Nump            int    //20
	Nums            int    //5
	Extra           string
	Title           string //字段描述
	OrdinalPosition string `json:"ordinalPosition" form:"ordinalPosition" gorm:"ordinalPosition"`
	IsSearchField   bool
	IsCreateField   bool
	IsPrimaryKey    bool
	IsIndex         bool
	DataMap         map[string]template.HTML `gorm:"-"`
	Jsontag         template.HTML `gorm:"-"`
	Gormtag         template.HTML `gorm:"-"`
	Formtag         template.HTML `gorm:"-"`
}

func NewColumn() *Column {
	return &Column{
		DataMap: make(map[string]template.HTML, 0),
	}
}
// 
func (ds *Column) Create() error {
	if ds.DbDataType == "PRI" {
		ds.IsPrimaryKey = true
	}
	return DbEngin.Create(ds).Error
}

func (ds *Column) BuildJsonTag() template.HTML {
	return template.HTML(`json:"` + ds.ColumnName + `"`)
}

func (ds *Column) BuildFormTag() template.HTML {
	return template.HTML(`form:"` + ds.ColumnName + `"`)
}

func (ds *Column) BuildGormTag() template.HTML {
	tmp := []string{}
	tmp = append(tmp, `gorm:"`+ds.DbColumnName+" "+ds.Extra+" "+ds.DbColumnType)
	if ds.IsPrimaryKey {
		tmp = append(tmp, "primaryKey")
	}
	if ds.IsIndex {
		tmp = append(tmp, "index")
	}
	tmp = append(tmp, "comment:"+ds.Title+"\"")
	return template.HTML(strings.Join(tmp, ";"))
}
func (ds *Column) BuildTag() template.HTML {
	return template.HTML(strings.Join([]string{"`", string(ds.BuildJsonTag()), string(ds.BuildFormTag()), string(ds.BuildGormTag()), "`"}, " "))
}

func (ds *Column) BuildLine() string {
	return strings.Join([]string{util.Ucfirst(ds.ColumnName), ds.DataType, string(ds.BuildTag())}, " ")
}

func (ds *Column) Search() (rows []Column, err error) {
	rows = make([]Column, 0)
	err = DbEngin.Model(ds).Find(&rows).Error
	return rows, err
}

func (ds *Column) Delete(ids []uint) (err error) {
	if len(ids) == 0 {
		return nil
	}
	err = DbEngin.Model(ds).Where("id in ?", ids).Delete(ds).Error
	return err
}
