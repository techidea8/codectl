package logic

import (
	"github.com/turingdance/infra/cond"
	"github.com/turingdance/infra/dbkit"
)

// 搜索
func Search[T any](model *T, wraper *cond.CondWraper) (result []T, total int64, err error) {
	return dbkit.Search(DbEngin, model, wraper)
}

// 创建一条记录
func Create[T any](model *T) (r *T, err error) {
	return dbkit.Create(DbEngin, model)
}

// 修改一条记录
func Update[T any](model *T, query interface{}, args ...interface{}) (r *T, err error) {
	return dbkit.Update(DbEngin, model, query, args...)
}

// 删除某条件
func Delete[T any](model *T, query interface{}, args ...interface{}) (effectrows int64, err error) {
	return dbkit.Delete(DbEngin, model, query, args...)
}

// 最先1条记录
func First[T any](model *T, wraper cond.CondWraper) (r *T, err error) {
	return dbkit.First(DbEngin, model, wraper)
}

// 最后一条记录
func Last[T any](model *T, wraper cond.CondWraper) (r *T, err error) {
	return dbkit.Last(DbEngin, model, wraper)

}
func Take[T any](model *T, wraper cond.CondWraper) (r *T, err error) {
	return dbkit.Take(DbEngin, model, wraper)
}

func TakeByPrimaryKey[T any](model *T) (r *T, err error) {
	return dbkit.TakeByPrimaryKey(DbEngin, model)
}
