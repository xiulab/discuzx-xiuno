package extension

import (
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/skiy/gfutils/llog"
	"strings"
)

// Forum 版块
type Forum struct {
}

// Parsing 解析
func (t *Forum) Parsing() (err error) {
	if cfg.GetBool("extension.forum.moderators") {
		return t.moderators()
	}
	return
}

// moderators 版主
func (t *Forum) moderators() (err error) {
	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxForumField := discuzPre + "forum_forumfield"

	fields := "fid,moderators"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxForumField).Where("moderators != ?", "").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.forum.name")
	xiunoUserTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")

	if err != nil {
		llog.Log.Debugf("表 %s 版主数据查询失败, %s", xiunoTable, err.Error())
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 无版主数据可以转换", xiunoTable)
		return nil
	}

	forumModArr := map[int][]string{}
	var moderArr g.SliceStr
	for _, u := range r.List() {
		moder := strings.Split(gconv.String(u["moderators"]), "	")

		// 有版主
		if len(moder) > 0 {
			forumModArr[gconv.Int(u["fid"])] = moder
			moderArr = append(moderArr, moder...)
		}
	}

	if len(moderArr) == 0 {
		llog.Log.Debugf("表 %s 无版主(moder)可以转换", xiunoTable)
		return nil
	}

	fields2 := "DISTINCT uid, username"
	r, err = database.GetXiunoDB().Table(xiunoUserTable).Where("username in (?)", moderArr).Fields(fields2).Select()
	if err != nil {
		llog.Log.Debugf("表 %s 版主用户名查询失败, %s", xiunoTable, err.Error())
	}

	if len(r) == 0 {
		llog.Log.Debug("表 %s 无版主用户可以转换", xiunoTable)
		return nil
	}

	uidArr := g.MapStrStr{}
	for _, u := range r.List() {
		uidArr[gconv.String(u["username"])] = gconv.String(u["uid"])
	}

	var modUIDArr g.SliceStr
	for fid, moders := range forumModArr {
		if len(moders) > 0 {
			for _, u := range moders {
				if uidArr[u] != "" {
					modUIDArr = append(modUIDArr, uidArr[u])
				}

				if len(modUIDArr) == 0 {
					continue
				}

				w := g.Map{
					"fid": fid,
				}

				d := g.Map{
					"moduids": strings.Join(modUIDArr, ","),
				}

				if _, err := database.GetXiunoDB().Table(xiunoTable).Data(d).Where(w).Update(); err != nil {
					return fmt.Errorf("表 %s 更新版主失败, %s", xiunoTable, err.Error())
				}
			}
		}
	}

	llog.Log.Infof("表 %s 更新版主成功", xiunoTable)
	return
}

// NewForum Forum init
func NewForum() *Forum {
	t := &Forum{}
	return t
}
