package extension

import (
	"database/sql"
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/skiy/gfutils/llog"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

// ThreadPost ThreadPost
type ThreadPost struct {
}

// NewThreadPost ThreadPost init
func NewThreadPost() *ThreadPost {
	t := &ThreadPost{}
	return t
}

// Parsing Parsing
func (t *ThreadPost) Parsing() (err error) {
	// 是否修正主题的 lastpid 和 lastuid
	if cfg.GetBool("extension.thread_post.fix_last") {
		if err := t.fixThreadLast(); err != nil {
			return err
		}
	}

	// 修正帖子内附件统计数量
	if cfg.GetBool("extension.thread_post.post_attach_total") {
		if err := t.postAttachTotal(); err != nil {
			return err
		}
	}

	// 修正主题内附件统计数量
	if cfg.GetBool("extension.thread_post.thread_attach_total") {
		if err := t.threadAttachTotal(); err != nil {
			return err
		}
	}
	return
}

// fixThreadLastGroup 是否修正主题的 lastpid 和 lastuid (分组导入)
func (t *ThreadPost) fixThreadLastGroup() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	xiunoThreadTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	xiunoDB := database.GetXiunoDB()

	fields := "max(pid) as max_pid"
	r, err := xiunoDB.Table(xiunoPostTable).Fields(fields).GroupBy("tid").Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 无数据可以转换 lastpid 和 lastuid", xiunoPostTable)
		return
	}

	// 分组导入
	var pidArr g.ArrayStr
	for _, u := range r.List() {
		pidArr = append(pidArr, gconv.String(u["max_pid"]))
	}

	// 获取最后一条帖子的 tid,uid,pid
	fields2 := "tid,pid,uid"
	res, err := xiunoDB.Table(xiunoPostTable).Where("pid in (?)", pidArr).Fields(fields2).Select()

	if err != nil {
		return err
	}

	if len(res) == 0 {
		llog.Log.Debugf("表 %s 找不到 lastpid 和 lastuid", xiunoPostTable)
		return
	}

	var count int64
	for _, u := range res.ToList() {
		w := g.Map{
			"tid": u["tid"],
		}

		d := g.Map{
			"lastpid": u["pid"],
			"lastuid": u["uid"],
		}

		var res2 sql.Result
		if res2, err = xiunoDB.Table(xiunoThreadTable).Data(d).Where(w).Update(); err != nil {
			return fmt.Errorf("表 %s 更新帖子的 lastuid 和 lastuid 失败, %s", xiunoThreadTable, err.Error())
		}

		c, _ := res2.RowsAffected()
		count += c
	}

	llog.Log.Infof("表 %s 更新帖子的 lastuid 和 lastuid 成功, 本次更新: %d 条数据, 耗时: %v", xiunoThreadTable, count, time.Since(start))
	return
}

// fixThreadLast 是否修正主题的 lastpid 和 lastuid
func (t *ThreadPost) fixThreadLast() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	xiunoThreadTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	xiunoDB := database.GetXiunoDB()

	fields := "max(pid) as max_pid"
	r, err := xiunoDB.Table(xiunoPostTable).Fields(fields).GroupBy("tid").Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 无数据可以转换 lastpid 和 lastuid", xiunoPostTable)
		return
	}

	var count int64
	// 获取最后一条帖子的 tid,uid,pid
	fields2 := "tid,pid,uid"
	for _, p := range r.List() {
		u, err := xiunoDB.Table(xiunoPostTable).Where(g.Map{"pid": p["max_pid"]}).Fields(fields2).One()
		if err != nil {
			return err
		}

		w := g.Map{
			"tid": u["tid"],
		}

		d := g.Map{
			"lastpid": u["pid"],
			"lastuid": u["uid"],
		}

		var res2 sql.Result
		if res2, err = xiunoDB.Table(xiunoThreadTable).Data(d).Where(w).Update(); err != nil {
			return fmt.Errorf("表 %s 更新帖子的 lastuid 和 lastuid 失败, %s", xiunoThreadTable, err.Error())
		}

		c, _ := res2.RowsAffected()
		count += c
	}

	llog.Log.Infof("表 %s 更新帖子的 lastuid 和 lastuid 成功, 本次更新: %d 条数据, 耗时: %v", xiunoThreadTable, count, time.Since(start))
	return
}

// threadAttachTotal 修正主题内附件统计数量
func (t *ThreadPost) threadAttachTotal() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	xiunoThreadTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	xiunoDB := database.GetXiunoDB()

	fields := "tid,files,images"
	r, err := xiunoDB.Table(xiunoPostTable).Where("isfirst = ?", 1).Fields(fields).Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 无 files 和 images 数据可以更新至 %s", xiunoPostTable, xiunoThreadTable)
		return
	}

	// ALTER TABLE `bbs_thread` CHANGE `images` `images` SMALLINT(6) NOT NULL DEFAULT '0';
	// ALTER TABLE `bbs_thread` CHANGE `files` `files` SMALLINT(6) NOT NULL DEFAULT '0';
	if _, err = xiunoDB.Exec("ALTER TABLE `bbs_thread` CHANGE `images` `images` SMALLINT(6) NOT NULL DEFAULT '0'"); err != nil {
		return fmt.Errorf("表 %s 转换 images 字段为 SMALLINT(6) 失败, %s", xiunoThreadTable, err.Error())
	}

	if _, err = xiunoDB.Exec("ALTER TABLE `bbs_thread` CHANGE `files` `files` SMALLINT(6) NOT NULL DEFAULT '0'"); err != nil {
		return fmt.Errorf("表 %s 转换 files 字段为 SMALLINT(6) 失败, %s", xiunoThreadTable, err.Error())
	}

	var count int64
	for _, u := range r.List() {
		w := g.Map{
			"tid": u["tid"],
		}

		d := g.Map{
			"files":  u["files"],
			"images": u["images"],
		}

		var res sql.Result
		if res, err = xiunoDB.Table(xiunoThreadTable).Data(d).Where(w).Update(); err != nil {
			return fmt.Errorf("表 %s 更新主题的附件数(files)和图片数(images)失败, %s", xiunoThreadTable, err.Error())
		}

		c, _ := res.RowsAffected()
		count += c
	}

	llog.Log.Infof("表 %s 更新主题的附件数(files)和图片数(images)成功, 本次更新: %d 条数据, 耗时: %v", xiunoThreadTable, count, time.Since(start))
	return
}

// postAttachTotal 修正帖子内附件统计数量
func (t *ThreadPost) postAttachTotal() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoAttachTable := xiunoPre + cfg.GetString("tables.xiuno.attach.name")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	xiunoDB := database.GetXiunoDB()

	fields := "count(*) as total, pid, isimage"
	r, err := xiunoDB.Table(xiunoAttachTable).Fields(fields).GroupBy("pid,isimage").Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		llog.Log.Debugf("表 %s 找不到附件数和图片数", xiunoAttachTable)
		return
	}

	var count int64
	var w, d g.Map
	for _, u := range r.List() {
		w = g.Map{
			"pid": u["pid"],
		}

		// 图片
		if u["isimage"] == 1 {
			d = g.Map{
				"images": u["total"],
			}
		} else { //非图片
			d = g.Map{
				"files": u["total"],
			}
		}

		var res sql.Result
		if res, err = xiunoDB.Table(xiunoPostTable).Data(d).Where(w).Update(); err != nil {
			return fmt.Errorf("表 %s 更新帖子的附件数(files)和图片数(images)失败, %s", xiunoPostTable, err.Error())
		}

		c, _ := res.RowsAffected()
		count += c
	}

	llog.Log.Infof("表 %s 更新帖子的附件数(files)和图片数(images)成功, 本次更新: %d 条数据, 耗时: %v", xiunoPostTable, count, time.Since(start))
	return
}
