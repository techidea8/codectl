package model

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Ds struct {
	Id       uint   `json:"id" form:"id" gorm:"primaryKey;int(11) auto_increment"`
	Database string `json:"database" form:"database"`
	Title    string `json:"title" form:"title"`
	Dbtype   string `json:"dbtype" form:"dbtype"`
	User     string `json:"user" form:"user"`
	Pass     string `json:"pass" form:"pass"`
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port" gorm:"default:3306"`
	Query    string `json:"query" form:"query"`
	Active   int    `json:"active" form:"avtive" gorm:"default:2"`
	Tpl      string `json:"tpl" form:"tpl"`
	Package  string `json:"package" form:"package"`
}

func (ds *Ds) Create() error {
	if ds.Title == "" {
		ds.Title = ds.Database
	}
	return DbEngin.Create(ds).Error
}
func (ds *Ds) Update() error {
	if ds.Id > 0 {
		return DbEngin.Where("id = ?", ds.Id).Updates(ds).Error
	} else if ds.Database != "" {
		return DbEngin.Where("database = ?", ds.Database).Updates(ds).Error
	} else {
		return errors.New("not suport ")
	}
}
func (ds *Ds) Remove() error {
	if ds.Id > 0 {
		return DbEngin.Where("id = ?", ds.Id).Delete(ds).Error
	} else if ds.Database != "" {
		return DbEngin.Where("database = ?", ds.Database).Delete(ds).Error
	} else {
		return errors.New("not suport ")
	}
}

func (ds *Ds) TriggleActive(status int) error {
	if status == 1 {
		DbEngin.Model(&Ds{}).Where("id > ?", 0).Update("active", 2)
		return DbEngin.Model(&Ds{Id: ds.Id}).Update("active", status).Error
	} else {
		return DbEngin.Model(&Ds{Id: ds.Id}).Update("active", status).Error
	}
}

func (ds *Ds) Search() (rows []Ds, err error) {
	rows = make([]Ds, 0)
	err = DbEngin.Model(ds).Find(&rows, ds).Error
	return rows, err
}

func (ds *Ds) BuildMysqlEngin() (db *gorm.DB, err error) {
	//"mysql", "user:password@/dbname"
	mysqlDs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", ds.User, ds.Pass, ds.Host, ds.Port, ds.Database, ds.Query)
	db, err = gorm.Open(mysql.Open(mysqlDs), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return
}

func (ds *Ds) BuildSqliteEngin() (db *gorm.DB, err error) {
	//"mysql", "user:password@/dbname"
	mysqlDs := fmt.Sprintf("%s?%s", ds.Database, ds.Query)
	db, err = gorm.Open(sqlite.Open(mysqlDs), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return
}

func (ds *Ds) BuildDbEngin() (db *gorm.DB, err error) {
	if ds.Dbtype == "mysql" {
		return ds.BuildMysqlEngin()
	} else if ds.Dbtype == "sqlite" || ds.Dbtype == "sqlite3" {
		return ds.BuildSqliteEngin()
	} else {
		return nil, errors.New("not suport")
	}
}
