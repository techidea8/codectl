package logic

import (
	"time"

	"github.com/techidea8/codectl/app/gen/model"
	"github.com/techidea8/codectl/infra/cond"
	"github.com/techidea8/codectl/infra/dbkit"
)

// 列表全部项目
func ListAllProject(wraper *cond.CondWraper) (result []model.Project, total int64, err error) {
	wraper.Pager.Pagesize = 1024 * 1024
	return dbkit.Search(DbEngin, &model.Project{}, wraper)
}

// 列表全部项目
func UseCurrentProject(prjId int32) (result *model.Project, err error) {
	return dbkit.Update(DbEngin, &model.Project{
		ID:        prjId,
		SortIndex: int32(time.Now().Unix()),
	}, "id = ?", prjId)
}

func UpdateProject(model *model.Project, query interface{}, args ...interface{}) (*model.Project, error) {
	return dbkit.Update(DbEngin, model, query, args...)
}
func DeleteProject(model *model.Project, query interface{}, args ...interface{}) (total int64, err error) {
	return dbkit.Delete(DbEngin, model, query, args...)
}

// 获得当前项目
func TakeProjectByPrimaryKey(instance *model.Project) (result *model.Project, err error) {
	return dbkit.Take(DbEngin, instance, cond.CondWraper{})
}

// 获得当前项目
func TakeCurrentProject() (result *model.Project, err error) {
	return dbkit.Take(DbEngin, &model.Project{}, cond.CondWraper{
		Order: cond.Order{
			Method: "desc",
			Field:  "sort_index",
		},
	})
}

func CreateProject(model *model.Project) (*model.Project, error) {
	return dbkit.Create(DbEngin, model)
}
