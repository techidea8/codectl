package dbkit

import (
	"github.com/techidea8/codectl/infra/cond"
	"gorm.io/gorm"
)

type Service struct {
	dbengin *gorm.DB
}

// 搜索
func (s *Service) Search(model interface{}, wraper *cond.CondWraper) (result interface{}, total int64, err error) {
	db := s.dbengin.Model(model)
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
	result = make([]interface{}, 0)
	err = db.Find(&result).Error
	return result, total, err
}

// 创建一条记录
func (s *Service) Create(model interface{}) (r interface{}, err error) {
	db := s.dbengin.Model(model)
	err = db.Create(model).Error
	return model, err
}

// 修改一条记录
func (s *Service) Update(model interface{}, query interface{}, args ...interface{}) (r interface{}, err error) {
	db := s.dbengin.Model(model)
	err = db.Where(query, args...).Updates(model).Error
	return model, err
}

// 删除某条件
func (s *Service) Delete(model interface{}, query interface{}, args ...interface{}) (effectrows int64, err error) {
	db := s.dbengin.Model(model)
	result := db.Where(model).Where(query, args...).Delete(model)
	return result.RowsAffected, result.Error
}

// 最先1条记录
func (s *Service) First(model interface{}, wraper cond.CondWraper) (r interface{}, err error) {
	db := s.dbengin
	for _, v := range wraper.Conds {
		sql, arg, err := v.Build()
		if err != nil {
			return nil, err
		}
		db = db.Where(sql, arg)
	}
	err = db.Limit(1).Where(model).First(model).Error
	return model, err
}

// 最后一条记录
func (s *Service) Last(model interface{}, wraper cond.CondWraper) (r interface{}, err error) {
	db := s.dbengin
	for _, v := range wraper.Conds {
		sql, arg, err := v.Build()
		if err != nil {
			return nil, err
		}
		db = db.Where(sql, arg)
	}
	err = db.Limit(1).Where(model).Last(model).Error
	return model, err

}
func (s *Service) Take(model interface{}, wraper cond.CondWraper) (r interface{}, err error) {
	db := s.dbengin
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

func (s *Service) TakeByPrimaryKey(model interface{}) (r interface{}, err error) {
	err = s.dbengin.Find(model).Limit(1).Error
	return model, err
}
