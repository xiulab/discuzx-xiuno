package extension

import (
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/skiy/gfutils/llog"
	"time"

	"github.com/gogf/gf/frame/g"
)

// Group 用户组
type Group struct {
}

// NewGroup Group init
func NewGroup() *Group {
	t := &Group{}
	return t
}

// Parsing 解析
func (t *Group) Parsing() (err error) {
	// 使用官方组
	if cfg.GetBool("tables.xiuno.group.official") {
		return t.official()
	}

	return nil
}

// official 官方组时,则需要将所有的用户组清理
func (t *Group) official() (err error) {
	xiunoPre := database.GetPrefix("xiuno")
	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")

	var start time.Time
	var count int64

	var d = g.Map{
		"gid": 101,
	}

	start = time.Now()
	r, err := database.GetXiunoDB().Table(xiunoTable).Data(d).Update()
	if err != nil {
		return fmt.Errorf("表 %s 重置[所有用户]的用户组 gid 为 101 失败, %s", xiunoTable, err.Error())
	}

	count, _ = r.RowsAffected()
	llog.Log.Infof("表 %s 重置[所有用户]的用户组 gid 为 101 成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start))

	// 不转换管理员 gid
	adminID := cfg.GetInt("extension.group.admin_id")
	if adminID <= 0 {
		return
	}

	d = g.Map{
		"gid": 1,
	}

	w := g.Map{
		"uid": adminID,
	}

	r, err = database.GetXiunoDB().Table(xiunoTable).Where(w).Data(d).Update()
	if err != nil {
		return fmt.Errorf("表 %s 重置 uid 为 %d 的用户组 gid 为 1 失败, %s", xiunoTable, adminID, err.Error())
	}

	count, _ = r.RowsAffected()
	llog.Log.Infof("表 %s 重置 uid 为 %d 的用户组为->管理员组(1) 成功", xiunoTable, adminID)
	return
}
