package biz

import (
	"fmt"

	"github.com/spf13/cobra" // 安装依赖 go get -u github.com/spf13/cobra/cobra
)

type sqlctrl struct {
}

func (s *sqlctrl) run(args []string) error {

	return nil
}

var sqlstr = `
-- 菜单 SQL
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('${functionName}', '${parentMenuId}', '1', '${businessName}', '${moduleName}/${businessName}/index', 1, 0, 'C', '0', '0', '${permissionPrefix}:list', '#', 'admin', sysdate(), '', null, '${functionName}菜单');
SELECT @parentId := LAST_INSERT_ID();
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('${functionName}查询', @parentId, '1',  '#', '', 1, 0, 'F', '0', '0', '${permissionPrefix}:query',        '#', 'admin', sysdate(), '', null, '');
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('${functionName}新增', @parentId, '2',  '#', '', 1, 0, 'F', '0', '0', '${permissionPrefix}:add',          '#', 'admin', sysdate(), '', null, '');
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('${functionName}修改', @parentId, '3',  '#', '', 1, 0, 'F', '0', '0', '${permissionPrefix}:edit',         '#', 'admin', sysdate(), '', null, '');
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('${functionName}删除', @parentId, '4',  '#', '', 1, 0, 'F', '0', '0', '${permissionPrefix}:remove',       '#', 'admin', sysdate(), '', null, '');
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('${functionName}导出', @parentId, '5',  '#', '', 1, 0, 'F', '0', '0', '${permissionPrefix}:export',       '#', 'admin', sysdate(), '', null, '');
`

// 子命令定义 运行方法 go run main.go version 编译后 ./hugo version
var sqlCmd = &cobra.Command{
	Use:   "sql", // Use这里定义的就是命令的名称
	Short: "sql runner",
	Long: `
sql run
    run a sql file 
`,
	Run: func(cmd *cobra.Command, args []string) { //这里是命令的执行方法
		fmt.Println("not implents!!")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行前执行
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行后执行
	},
	// 还有其他钩子函数
}

func init() {
	rootCmd.AddCommand(sqlCmd)
	sqlCmd.Flags().Int32VarP(&prjId, "prj", "p", 0, "id of project ")
}
