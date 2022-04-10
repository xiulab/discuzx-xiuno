package controllers

import (
	"database/sql"
	"discuzx-xiuno/app/libraries/common"
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
	"regexp"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/util/gconv"
)

// User User
type User struct {
}

// ToConvert user ToConvert
func (t *User) ToConvert() (err error) {
	cfg := lcfg.Get()

	if cfg.GetString("database.discuz.0.host") == cfg.GetString("database.uc.0.host") &&
		cfg.GetString("database.discuz.0.port") == cfg.GetString("database.uc.0.port") &&
		cfg.GetString("database.discuz.0.user") == cfg.GetString("database.uc.0.user") &&
		cfg.GetString("database.discuz.0.name") == cfg.GetString("database.uc.0.name") {

		llog.Log.Debugf("Discuz & UCenter 是同一个数据库")

		return t.sameUCenter()
	}

	llog.Log.Debugf("Discuz & UCenter 不是同一个数据库")
	return t.otherUCenter()
}

// sameUCenter UCenter 与 Discuz!X 同一个库
func (t *User) sameUCenter() (err error) {
	start := time.Now()

	cfg := lcfg.Get()

	ucPre, discuzPre, xiunoPre := database.GetPrefix("uc"), database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	ucMemberTable := ucPre + "members"
	dxMemberTable := discuzPre + "common_member"
	dxMemberStatusTable := discuzPre + "common_member_status"

	fields := "m.uid,m.groupid,m.email,m.username,m.credits,m.regdate,s.regip,s.lastip,s.lastvisit,u.password,u.salt"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxMemberTable+" m").LeftJoin(dxMemberStatusTable+" s", "s.uid = m.uid").LeftJoin(ucMemberTable+" u", "u.uid = m.uid").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")
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

	xiunoDB := database.GetXiunoDB()
	if _, err = xiunoDB.Exec("TRUNCATE " + xiunoTable); err != nil {
		return fmt.Errorf("清空数据表(%s)失败, %s", xiunoTable, err.Error())
	}

	multipleEmailFlag := cfg.GetInt("tables.xiuno.user.multiple_email")
	// 清除索引方式
	if multipleEmailFlag == 1 {
		if _, err := xiunoDB.Exec("ALTER TABLE " + xiunoTable + " DROP INDEX email"); err != nil {
			llog.Log.Noticef("表 %s 清除 email 唯一索引失败, %s", xiunoTable, err.Error())
		}
	}

	var count int64
	batch := cfg.GetInt("tables.xiuno.user.batch")

	// 修改 email 方式去重, 只能一条条地加
	if multipleEmailFlag == 2 {
		batch = 1
	}

	emailMul := `for key 'email'`
	regEmailMul, _ := regexp.Compile(emailMul)

	emailLong := `Data too long for column 'email'`
	regEmailLong, _ := regexp.Compile(emailLong)

	dataList := gdb.List{}
	for _, u := range r.List() {
		password := gconv.String(u["password"])
		if password == "" {
			password = "mustResetPassword"
		}

		salt := gconv.String(u["salt"])
		if salt == "" {
			salt = common.GetRandomString("numeric", 6)
		}

		uid := gconv.Int(u["uid"])
		email := gconv.String(u["email"])

		d := gdb.Map{
			"uid":         uid,
			"gid":         common.FixGID(u["groupid"]),
			"email":       email,
			"username":    u["username"],
			"password":    password,
			"salt":        salt,
			"credits":     u["credits"],
			"create_ip":   common.IP2Long(gconv.String(u["regip"])),
			"create_date": gconv.Int(u["regdate"]),
			"login_ip":    common.IP2Long(gconv.String(u["lastip"])),
			"login_date":  gconv.Int(u["lastvisit"]),
		}

		// 批量插入数量
		if batch > 1 {
			dataList = append(dataList, d)
		} else {
			res, err := xiunoDB.Insert(xiunoTable, d)
			if err != nil {
				var isDo bool

				// email 过长
				if isNil := regEmailLong.FindString(err.Error()); isNil != "" {
					if _, err := xiunoDB.Exec("ALTER TABLE `" + xiunoTable + "` CHANGE `email` `email` CHAR(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '邮箱'"); err != nil {
						llog.Log.Noticef("表 %s 更改 email 长度为 255 失败, %s", xiunoTable, err.Error())
					}

					res, err = xiunoDB.Insert(xiunoTable, d)
					isDo = true
				}

				// email 重复, 修改 email 方式, 则再重新提交一次
				if !isDo && multipleEmailFlag == 2 {
					// 错误是 email 重复
					if isNil := regEmailMul.FindString(err.Error()); isNil != "" {
						d["email"] = fmt.Sprintf("%d_%s", uid, email)
						res, err = xiunoDB.Insert(xiunoTable, d)
						isDo = true
					}
				}

				if err != nil {
					return fmt.Errorf("表 %s 数据导入失败, %s", xiunoTable, err.Error())
				}
			}

			c, _ := res.RowsAffected()
			count += c
		}
	}

	if len(dataList) > 0 {
		res, err := xiunoDB.BatchInsert(xiunoTable, dataList, batch)
		if err != nil {
			return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
		}

		count, _ = res.RowsAffected()
	}

	llog.Log.Infof("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start))
	return
}

// otherUCenter UCenter 与 Discuz!X 不同一个库
func (t *User) otherUCenter() (err error) {
	start := time.Now()

	cfg := lcfg.Get()

	ucPre, discuzPre, xiunoPre := database.GetPrefix("uc"), database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	ucMemberTable := ucPre + "members"
	dxMemberTable := discuzPre + "common_member"
	dxMemberStatusTable := discuzPre + "common_member_status"

	fields := "m.uid,m.groupid,m.email,m.username,m.credits,m.regdate,s.regip,s.lastip,s.lastvisit"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxMemberTable+" m").LeftJoin(dxMemberStatusTable+" s", "s.uid = m.uid").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")
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

	xiunoDB := database.GetXiunoDB()
	if _, err = xiunoDB.Exec("TRUNCATE " + xiunoTable); err != nil {
		return fmt.Errorf("表 %s 清空失败, %s", xiunoTable, err.Error())
	}

	multipleEmailFlag := cfg.GetInt("tables.xiuno.user.multiple_email")
	// 清除索引方式
	if multipleEmailFlag == 1 {
		if _, err := xiunoDB.Exec("ALTER TABLE " + xiunoTable + " DROP INDEX email"); err != nil {
			llog.Log.Noticef("表 %s 清除 email 唯一索引失败, %s", xiunoTable, err.Error())
		}
	}

	var count int64
	fields2 := "password,salt"

	emailMul := `for key 'email'`
	regEmailMul, _ := regexp.Compile(emailMul)

	emailLong := `Data too long for column 'email'`
	regEmailLong, _ := regexp.Compile(emailLong)

	for _, u := range r.List() {
		password := "mustResetPassword" // 默认密码
		salt := ""                      // 盐值

		// 查询密码
		w2 := gdb.Map{
			"uid": u["uid"],
		}
		r2, err := database.GetUcDB().Table(ucMemberTable).Where(w2).Fields(fields2).One()
		// 无错误,且有数据
		if err == nil && r2 != nil {
			password = gconv.String(r2["password"])
			salt = gconv.String(r2["salt"])
		}

		if salt == "" {
			salt = common.GetRandomString("numeric", 6)
		}

		uid := gconv.Int(u["uid"])
		email := gconv.String(u["email"])

		d := gdb.Map{
			"uid":         uid,
			"gid":         common.FixGID(u["groupid"]),
			"email":       email,
			"username":    u["username"],
			"password":    password,
			"salt":        salt,
			"credits":     u["credits"],
			"create_ip":   common.IP2Long(gconv.String(u["regip"])),
			"create_date": gconv.Int(u["regdate"]),
			"login_ip":    common.IP2Long(gconv.String(u["lastip"])),
			"login_date":  gconv.Int(u["lastvisit"]),
		}

		res, err := xiunoDB.Insert(xiunoTable, d)
		if err != nil {
			var isDo bool

			// email 过长
			if isNil := regEmailLong.FindString(err.Error()); isNil != "" {
				if _, err := xiunoDB.Exec("ALTER TABLE `" + xiunoTable + "` CHANGE `email` `email` CHAR(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '邮箱'"); err != nil {
					llog.Log.Noticef("表 %s 更改 email 长度为 255 失败, %s", xiunoTable, err.Error())
				}

				res, err = xiunoDB.Insert(xiunoTable, d)
				isDo = true
			}

			// email 重复, 修改 email 方式, 则再重新提交一次
			if !isDo && multipleEmailFlag == 2 {
				// 错误是 email 重复
				if isNil := regEmailMul.FindString(err.Error()); isNil != "" {
					d["email"] = fmt.Sprintf("%d_%s", uid, email)
					res, err = xiunoDB.Insert(xiunoTable, d)
					isDo = true
				}
			}

			if err != nil {
				return fmt.Errorf("表 %s 数据导入失败, %s", xiunoTable, err.Error())
			}
		}

		c, _ := res.RowsAffected()
		count += c
	}

	llog.Log.Info("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start))
	return
}

// NewUser User init
func NewUser() *User {
	t := &User{}
	return t
}
