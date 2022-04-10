package controllers

import (
	"database/sql"
	"discuzx-xiuno/app/libraries/common"
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/util/gconv"
)

// Group Group
type Group struct {
}

// ToConvert ToConvert
func (t *Group) ToConvert() (err error) {
	cfg := lcfg.Get()

	officialSQL := `
TRUNCATE TABLE bbs_group;
INSERT INTO bbs_group (gid, name, creditsfrom, creditsto, allowread, allowthread, allowpost, allowattach, allowdown, allowtop, allowupdate, allowdelete, allowmove, allowbanuser, allowdeleteuser, allowviewip) VALUES
(0, '游客组', 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0),
(1, '管理员组', 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
(2, '超级版主组', 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1),
(4, '版主组', 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1),
(5, '实习版主组', 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0),
(6, '待验证用户组', 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0),
(7, '禁止用户组', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
(101, '一级用户组', 0, 50, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0),
(102, '二级用户组', 50, 200, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0),
(103, '三级用户组', 200, 1000, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0),
(104, '四级用户组', 1000, 10000, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0),
(105, '五级用户组', 10000, 10000000, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0);
`

	xiunoDB := database.GetXiunoDB()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.group.name")

	// 重置官方用户组
	if _, err = xiunoDB.Exec(officialSQL); err != nil {
		return fmt.Errorf("重置官方用户组(%s)失败, %s", xiunoTable, err.Error())
	}

	// 使用 XiunoBBS 官方用户组, 则不转换
	if cfg.GetBool("tables.xiuno.group.official") {
		return
	}

	start := time.Now()

	dxGroupTable := discuzPre + "common_usergroup"
	dxGroupField := discuzPre + "common_usergroup_field"
	dxAdminGroup := discuzPre + "common_admingroup"

	fields := "u.groupid,u.grouptitle,u.type,u.creditslower,u.creditshigher,u.allowvisit,f.allowpost,f.allowreply,f.allowpostattach,f.allowgetattach,a.allowstickthread,a.alloweditpost,a.allowdelpost,a.allowmovethread,a.allowbanvisituser,a.allowviewip"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxGroupTable+" u").LeftJoin(dxGroupField+" f", "f.groupid = u.groupid").LeftJoin(dxAdminGroup+" a", "a.admingid = u.groupid").Fields(fields).Select()

	if err != nil {
		if err == sql.ErrNoRows {
			llog.Log.Debugf("表 %s 无数据可以转换", xiunoTable)
			return nil
		}

		llog.Log.Debugf("表 %s 数据查询失败, %s", xiunoTable, err.Error())
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 无数据可以转换", xiunoTable)
		return nil
	}

	// 删除 101 以外的用户组
	if _, err = xiunoDB.Table(xiunoTable).Where("gid > ", 101).Delete(); err != nil {
		return fmt.Errorf("清空数据表(%s)失败, %s", xiunoTable, err.Error())
	}

	var count int64
	dataList := gdb.List{}
	for _, u := range r.List() {
		allowtop := gconv.Int(u["allowstickthread"])

		gid := common.FixGID(u["groupid"])

		// 小于默认用户组则全跳过
		if gid < 102 {
			continue
		}

		d := gdb.Map{
			"gid":          gid,
			"name":         u["grouptitle"],
			"creditsfrom":  u["creditshigher"],
			"creditsto":    u["creditslower"],
			"allowread":    u["allowvisit"],
			"allowthread":  u["allowpost"],
			"allowpost":    u["allowreply"],
			"allowattach":  u["allowpostattach"],
			"allowdown":    u["allowgetattach"],
			"allowtop":     allowtop,
			"allowupdate":  u["alloweditpost"],
			"allowdelete":  u["allowdelpost"],
			"allowmove":    u["allowmovethread"],
			"allowbanuser": u["allowbanvisituser"],
			"allowviewip":  u["allowviewip"],
		}

		// 普通会员
		if gconv.String(u["type"]) == "member" {
			d["allowtop"] = 0
			d["allowupdate"] = "0"
			d["allowdelete"] = "0"
			d["allowmove"] = "0"
			d["allowbanuser"] = "0"
			d["allowviewip"] = "0"
		}

		// 允许置顶,则值全为 1
		if gconv.Int(d["allowtop"]) > 0 {
			d["allowtop"] = 1
		}

		dataList = append(dataList, d)
	}

	if len(dataList) > 0 {
		res, err := xiunoDB.BatchInsert(xiunoTable, dataList, 100)
		if err != nil {
			return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
		}
		count, _ = res.RowsAffected()
	}

	llog.Log.Infof("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start))
	return
}

// NewGroup Group init
func NewGroup() *Group {
	t := &Group{}
	return t
}
