package extension

import (
	"database/sql"
	"discuzx-xiuno/app/libraries/common"
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// User User
type User struct {
}

// NewUser User init
func NewUser() *User {
	t := &User{}
	return t
}

// Parsing Parsing
func (t *User) Parsing() (err error) {
	// 是否修正主题 (post.first 全部为0,最第一条变更为主题)
	if cfg.GetBool("extension.user.fix_thread") {
		if err := t.fixThread(); err != nil {
			return err
		}
	}

	// 修正 gid 为 101 的用户及用户组
	if cfg.GetBool("extension.user.normal_user") {
		if err := t.normalUser(); err != nil {
			return err
		}
	}

	// 修正用户主题和帖子统计
	if cfg.GetBool("extension.user.total") {
		err := t.threadPostStat()
		return err
	}
	return
}

// fixThread 是否修正主题 (post.first 全部为0,最第一条变更为主题)
func (t *User) fixThread() (err error) {
	start := time.Now()

	cfg := lcfg.Get()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxName := cfg.GetString("database.discuz.0.name")
	xnName := cfg.GetString("database.xiuno.0.name")

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")

	dxPostTable := fmt.Sprintf("%s.%sforum_post", dxName, discuzPre)
	dxThreadTable := fmt.Sprintf("%s.%sforum_thread", dxName, discuzPre)
	xnThreadTable := fmt.Sprintf("%s.%s", xnName, xiunoTable)

	xiunoDB := database.GetXiunoDB()

	// SELECT x.tid xtid, d.tid FROM gearer.pre_forum_thread d LEFT JOIN xn.bbs_thread x ON x.tid = d.tid WHERE x.tid IS NULL
	// SELECT x.tid xtid, d.tid,p.pid FROM gearer.pre_forum_thread d LEFT JOIN xn.bbs_thread x ON x.tid = d.tid LEFT JOIN xn.bbs_post p ON p.tid = d.tid WHERE x.tid IS NULL GROUP BY tid
	fields := "t.fid,t.tid,t.displayorder,t.subject,t.dateline,t.lastpost,t.views,t.replies,t.closed,p.useip,p.pid,p.authorid"
	var r gdb.Result
	r, err = xiunoDB.Table(dxThreadTable+" t").LeftJoin(xnThreadTable+" x", "x.tid = t.tid").LeftJoin(dxPostTable+" p", "p.tid = t.tid").Where("x.tid IS NULL AND p.pid IS NOT NULL").GroupBy("tid").OrderBy("t.tid ASC").Fields(fields).Select()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		llog.Log.Debugf("表 %s 数据查询失败, %s", xnThreadTable, err.Error())
	}

	if len(r) == 0 {
		return nil
	}

	var count int64
	countMax := len(r.List())

	for _, u := range r.List() {
		tid := gconv.Int(u["tid"])
		userip := common.IP2Long(gconv.String(u["useip"]))

		d := gdb.Map{
			"fid":         u["fid"],
			"tid":         tid,
			"top":         u["displayorder"],
			"uid":         u["authorid"],
			"userip":      userip,
			"subject":     u["subject"],
			"create_date": u["dateline"],
			"last_date":   u["lastpost"],
			"views":       u["views"],
			"posts":       u["replies"],
			"closed":      u["closed"],
			"firstpid":    u["pid"],
		}

		res, err := xiunoDB.Insert(xiunoTable, d)
		if err != nil {
			return fmt.Errorf("[修正]表 %s 数据插入失败, %s", xiunoTable, err.Error())
		}
		c, _ := res.RowsAffected()
		count += c
	}

	llog.Log.Infof("[修正]表 %s 数据导入成功, 本次导入: %d/%d 条数据, 耗时: %v", xiunoTable, count, countMax, time.Since(start))
	return
}

// normalUser 修正 gid 为 101 的用户及用户组
func (t *User) normalUser() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoGroupTable := xiunoPre + cfg.GetString("tables.xiuno.group.name")
	xiunoUserTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")

	fields := "gid,name"
	r, err := database.GetXiunoDB().Table(xiunoGroupTable).Where("creditsfrom = ? AND creditsto > ?", 0, 0).Fields(fields).OrderBy("gid ASC").One()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 无用户组可以转换", xiunoGroupTable)
		return
	}

	var d, w g.Map

	d = g.Map{
		"gid": 101,
	}

	w = g.Map{
		"gid": r["gid"],
	}

	_, err = database.GetXiunoDB().Table(xiunoGroupTable).Where(w).Data(d).Update()
	if err != nil {
		return fmt.Errorf("表 %s 原 “%v” 组(%v) 转换为普通用户组 gid 为 101 失败, %s", xiunoGroupTable, r["name"], r["gid"], err.Error())
	}
	llog.Log.Infof("表 %s 原 “%v” 组(%v) 转换为普通用户组 gid 为 101 成功", xiunoGroupTable, r["name"], r["gid"])

	res, err := database.GetXiunoDB().Table(xiunoUserTable).Where(w).Data(d).Update()
	if err != nil {
		return fmt.Errorf("表 %s 原 “%v” 组(%v) 的用户转换为普通用户组 gid 为 101 失败, %s", xiunoGroupTable, r["name"], r["gid"], err.Error())
	}
	count, _ := res.RowsAffected()
	llog.Log.Infof("表 %s 原 “%v” 组(%v)的用户转换为普通用户组 gid 为 101 成功, 本次更新: %d 条数据, 耗时: %v", xiunoGroupTable, r["name"], r["gid"], count, time.Since(start))

	return
}

// threadPostStat 修正用户主题和帖子数量, 帖子包含主题和回复
func (t *User) threadPostStat() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoUserTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")
	xiunoThreadTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")

	xiunoDB := database.GetXiunoDB()

	fields := "uid"
	var r gdb.Result
	r, err = xiunoDB.Table(xiunoUserTable).Fields(fields).Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 无用户可以转换主题和帖子数量", xiunoUserTable)
		return
	}

	var count int64
	for _, u := range r.List() {
		w := g.Map{
			"uid": u["uid"],
		}
		posts, err := database.GetXiunoDB().Table(xiunoPostTable).Where(w).Fields("tid").Count()
		if err != nil {
			posts = 0
		}

		threads, err := database.GetXiunoDB().Table(xiunoThreadTable).Where(w).Fields("tid").Count()
		if err != nil {
			threads = 0
		}

		d := g.Map{
			"threads": threads,
			"posts":   posts,
		}

		res, err := xiunoDB.Table(xiunoUserTable).Data(d).Where(w).Update()
		if err != nil {
			return fmt.Errorf("表 %s 用户帖子统计更新失败, %s", xiunoUserTable, err.Error())
		}
		c, _ := res.RowsAffected()
		count += c
	}

	llog.Log.Infof("表 %s 用户帖子统计, 本次更新: %d 条数据, 耗时: %v", xiunoUserTable, count, time.Since(start))
	return
}
