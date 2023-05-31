package model

import (
	"errors"

	"gorm.io/gorm"
	"techidea8.com/codectl/internal/util"
)

type Table struct {
	Id        uint     `json:"id" form:"id" gorm:"primaryKey;int(11) auto_increment"`
	Database  string   `json:"database" form:"database"`
	TableName string   `json:"tableName" form:"tableName"`
	Title     string   `json:"title" form:"title"`
	Module    string   `json:"module" form:"module"`
	Package   string   `json:"package" form:"package"`
	Tpldir    string   `json:"tpldir" form:"tpldir"`
	Savedir   string   `json:"savedir" form:"savedir"`
	Method    []Method `json:"method" form:"method" gorm:"serializer:json"`
	Columns   []Column `json:"columns" form:"column" gorm:"-"`
}

func (instance *Table) BeforeCreate(tx *gorm.DB) error {
	if instance.Title == "" {
		instance.Title = instance.Module
	}
	return nil
}

func (instance *Table) CreateDs() Ds {
	tmp := Ds{
		Database: instance.Database,
	}
	DbEngin.Where("database = ?", instance.Database).First(&tmp)
	return tmp
}
func (instance *Table) Create() error {
	if len(instance.Method) == 0 {
		instance.Method = make([]Method, 0)
		instance.Method = append(instance.Method, DefaultMethods...)
	}
	if instance.TableName == "" {
		return errors.New("please set tableName use -t xxx")
	}
	if instance.Module == "" {
		instance.Module = util.UnderlineToCamelCase(instance.TableName)
	}
	tmp := Table{
		Database:  instance.Database,
		TableName: instance.TableName,
	}
	DbEngin.First(&tmp)
	if tmp.Id > 0 {
		instance.Id = tmp.Id
	}
	return DbEngin.Save(instance).Error
}

func (instance *Table) Search() (rows []Table, err error) {
	rows = make([]Table, 0)
	err = DbEngin.Model(instance).Find(&rows, instance).Error
	return rows, err
}
