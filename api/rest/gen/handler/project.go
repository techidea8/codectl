package handler

import (
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/techidea8/codectl/infra/cond"
	"github.com/techidea8/codectl/infra/restkit"
	"github.com/techidea8/codectl/infra/wraper"
)

type Project struct{}

func init() {
	restkit.Register(&Project{})
}

// 创建项目
func (prj *Project) Create(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Project{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	instance, err = logic.CreateProject(instance)
	return wraper.OkData(instance).WithMsg("项目创建成功"), err
}

func (prj *Project) Update(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Project{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	instance, err = logic.UpdateProject(instance, "id = ?", instance.ID)
	return wraper.OkData(instance).WithMsg("项目修改成功"), err
}

func (prj *Project) List(ctx restkit.Context) (r *wraper.Response, err error) {
	prjs, total, err := logic.ListAllProject(&cond.CondWraper{
		Pager: cond.Pager{
			Pagesize: 1024,
		},
		Order: cond.Order{
			Field:  "sort_index",
			Method: "desc",
		},
	})
	return wraper.OkData(prjs).WithTotal(total), err
}

func (prj *Project) Search(ctx restkit.Context) (r *wraper.Response, err error) {
	conds := cond.NewCondWrapper()
	if err = ctx.Bind(conds); err != nil {
		return
	}
	prjs, total, err := logic.Search(&model.Project{}, conds)
	return wraper.OkData(prjs).WithTotal(total), err
}

func (prj *Project) Delete(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Project{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	total, err := logic.DeleteProject(instance, "id = ?", instance.ID)
	return wraper.OkData(total).WithMsg("项目删除成功"), err
}

func (prj *Project) Meta(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Table{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	metas, err := logic.MetaFromModuleName(instance.Module)
	return wraper.OkData(metas), err
}

func (prj *Project) Clone(ctx restkit.Context) (r *wraper.Response, err error) {
	instance := &model.Project{}
	err = ctx.Bind(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	result, err := logic.TakeProjectByPrimaryKey(instance)
	if err != nil {
		return wraper.Error(err), err
	}
	result.ID = 0
	result, err = logic.CreateProject(result)
	return wraper.OkData(result).WithMsg("项目克隆成功"), err
}
