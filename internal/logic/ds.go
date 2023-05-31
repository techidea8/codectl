package logic

import (
	"encoding/json"

	"techidea8.com/codectl/internal/model"
	"techidea8.com/codectl/internal/util"
)

func TakeDs(dnsId any) (ds model.Ds, err error) {
	ds = model.Ds{}
	if dsId, ok := dnsId.(uint); ok {
		err = model.DbEngin.First(&ds, "id = ?", dsId).Error
	}
	if dbname, ok := dnsId.(string); ok {
		err = model.DbEngin.First(&ds, "database = ?", dbname).Error
	}
	return ds, err
}

func ListDs() (ds []model.Ds, err error) {
	ds = make([]model.Ds, 0)
	err = model.DbEngin.Model(&model.Ds{}).Find(&ds).Error
	return ds, err
}

func ListDns(in *model.Ds) error {
	ds, err := ListDs()
	if err != nil {
		return err
	}
	arr, err := json.Marshal(ds)
	if err != nil {
		return err
	}

	return util.PrintTableWithByteArray(arr)
}
func ActiveDns(ds *model.Ds) (err error) {
	err = ds.TriggleActive(1)
	ListDns(ds)
	return err
}

func UpdateDns(ds *model.Ds) (err error) {
	err = ds.Update()
	ListDns(ds)
	return err
}
func CreateDns(ds *model.Ds) (err error) {
	err = ds.Create()
	ListDns(ds)
	return err
}
func RemoveDns(ds *model.Ds) (err error) {
	err = ds.Remove()
	ListDns(ds)
	return
}
