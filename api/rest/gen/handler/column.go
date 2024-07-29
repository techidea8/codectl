package handler

import (
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/techidea8/codectl/infra/cond"
	"github.com/techidea8/codectl/infra/restkit"
	"github.com/techidea8/codectl/infra/wraper"
)

type column struct{}

// 创建字段
func (prj *column) Create(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Column{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	instance, err = logic.CreateColumn(instance)
	return wraper.OkData(instance).WithMsg("模块创建成功"), err
}

// 修改某一个字段
func (prj *column) Update(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Column{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	instance, err = logic.UpdateColumn(instance, "id = ?", instance.ID)
	return wraper.OkData(instance).WithMsg("模块修改成功"), err
}

// 列举全部字段
func (prj *column) List(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Column{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	prjs, total, err := logic.ListAllColumn(&cond.CondWraper{
		Conds: []cond.Cond{
			cond.Cond{
				Field: "table_id",
				Op:    cond.OPEQ,
				Value: instance.TableID,
			},
		},
		Pager: cond.Pager{
			Pagesize: 1024,
		},
	})
	return wraper.OkData(prjs).WithTotal(total), err
}

// 删除字段
func (prj *column) Delete(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Column{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	total, err := logic.DeleteColumn(instance, "id = ?", instance.ID)
	return wraper.OkData(total).WithMsg("模块删除成功"), err
}
