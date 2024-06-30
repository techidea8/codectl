package dbkit

import (
	"github.com/techidea8/codectl/infra/cond"
	"gorm.io/gorm"
)

// 搜索
func Search[T any](dbengin *gorm.DB, model *T, wraper *cond.CondWraper) (result []T, total int64, err error) {
	db := dbengin.Model(model)
	for _, v := range wraper.Conds {
		sql, arg, err := v.Build()
		if err != nil {
			return nil, 0, err
		}
		db = db.Where(sql, arg)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	db = db.Limit(wraper.Pager.Limit()).Offset(wraper.Pager.Offset())
	if order, err := wraper.Order.Build(); err == nil {
		db = db.Order(order)
	}
	result = make([]T, 0)
	err = db.Find(&result).Error
	return result, total, err
}

// 创建一条记录
func Create[T any](dbengin *gorm.DB, model *T) (r *T, err error) {
	db := dbengin.Model(model)
	err = db.Create(model).Error
	return model, err
}

// 修改一条记录
func Update[T any](dbengin *gorm.DB, model *T, query interface{}, args ...interface{}) (r *T, err error) {
	db := dbengin.Model(model)
	err = db.Where(query, args...).Updates(model).Error
	return model, err
}

// 删除某条件
func Delete[T any](dbengin *gorm.DB, model *T, query interface{}, args ...interface{}) (effectrows int64, err error) {
	db := dbengin.Model(model)
	result := db.Where(model).Where(query, args...).Delete(model)
	return result.RowsAffected, result.Error
}

// 最先1条记录
func First[T any](dbengin *gorm.DB, model *T, wraper cond.CondWraper) (r *T, err error) {
	db := dbengin
	for _, v := range wraper.Conds {
		sql, arg, err := v.Build()
		if err != nil {
			return nil, err
		}
		db = db.Where(sql, arg)
	}
	err = dbengin.Limit(1).Where(model).First(model).Error
	return model, err
}

// 最后一条记录
func Last[T any](dbengin *gorm.DB, model *T, wraper cond.CondWraper) (r *T, err error) {
	db := dbengin
	for _, v := range wraper.Conds {
		sql, arg, err := v.Build()
		if err != nil {
			return nil, err
		}
		db = db.Where(sql, arg)
	}
	err = dbengin.Limit(1).Where(model).Last(model).Error
	return model, err

}
func Take[T any](dbengin *gorm.DB, model *T, wraper cond.CondWraper) (r *T, err error) {
	db := dbengin
	for _, v := range wraper.Conds {
		sql, arg, err := v.Build()
		if err != nil {
			return nil, err
		}
		db = db.Where(sql, arg)
	}
	if order, err := wraper.Order.Build(); err == nil && order != "" {
		db = db.Order(order)
	}
	err = db.Limit(1).Where(model).Find(model).Error
	return model, err
}

func TakeByPrimaryKey[T any](dbengin *gorm.DB, model *T) (r *T, err error) {
	err = dbengin.Find(model).Limit(1).Error
	return model, err
}
