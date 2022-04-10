package controllers

import (
	"database/sql"
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/util/gconv"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
	"path"
	"strings"
	"time"
)

// Attach 附件
type Attach struct {
}

// ToConvert 附件转换
func (t *Attach) ToConvert() (err error) {
	start := time.Now()

	cfg := lcfg.Get()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxAttachmentTable := discuzPre + "forum_attachment"

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.attach.name")

	xiunoDB := database.GetXiunoDB()
	if _, err = xiunoDB.Exec("TRUNCATE " + xiunoTable); err != nil {
		return fmt.Errorf("清空数据表(%s)失败, %s", xiunoTable, err.Error())
	}

	var r gdb.Result
	var count int64
	var failureData []string

	fields := "aid,tid,pid,uid,filesize,width,attachment,filename,dateline,description,isimage"
	batch := cfg.GetInt("tables.xiuno.attach.batch")

	// forum_attachment_0 ~ forum_attachment_9
	for i := 0; i < 10; i++ {
		tbname := fmt.Sprintf("%s_%d", dxAttachmentTable, i)
		r, err = database.GetDiscuzDB().Table(tbname).Fields(fields).Select()

		if err != nil {
			if err == sql.ErrNoRows {
				llog.Log.Debugf("表 %s 无数据可以转换", xiunoTable)
				return nil
			}
			return fmt.Errorf("表 %s 数据查询失败, %s", tbname, err.Error())
		}

		if len(r) == 0 {
			llog.Log.Debugf("表 %s 无数据可以转换(%s)", xiunoTable, tbname)
			continue
		}

		dataList := gdb.List{}
		for _, u := range r.List() {
			filetype := t.FileExt(gconv.String(u["filename"]))

			isimage := gconv.Int(u["isimage"])
			if isimage != 1 {
				isimage = 0
			}

			d := gdb.Map{
				"aid":         u["aid"],
				"tid":         u["tid"],
				"pid":         u["pid"],
				"uid":         u["uid"],
				"filesize":    u["filesize"],
				"width":       u["width"],
				"filename":    u["attachment"],
				"orgfilename": u["filename"],
				"filetype":    filetype,
				"create_date": u["dateline"],
				"comment":     u["description"],
				"isimage":     isimage,
			}

			// 批量插入数量
			if batch > 1 {
				dataList = append(dataList, d)
			} else {
				if res, err := xiunoDB.Insert(xiunoTable, d); err != nil {
					//return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
					llog.Log.Noticef("表 %s 数据插入失败, %s", xiunoTable, err.Error())
					failureData = append(failureData, fmt.Sprintf("%s(aid:%v)", tbname, u["aid"]))
				} else {
					c, _ := res.RowsAffected()
					count += c
				}
			}
		}

		if len(dataList) > 0 {
			// 批量插入
			res, err := xiunoDB.BatchInsert(xiunoTable, dataList, batch)
			if err != nil {
				return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
			}
			c, _ := res.RowsAffected()
			count += c
		}
	}

	llog.Log.Infof("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start))

	if len(failureData) > 0 {
		llog.Log.Noticef("导入失败的数据: %v", failureData)
	}
	return nil
}

// FileExt 获取文件后缀
func (t *Attach) FileExt(filename string) string {
	fileSuffix := path.Ext(filename)
	suffix := strings.Replace(fileSuffix, ".", "", -1)

	fileext := "other"

	if strings.EqualFold("png", suffix) || strings.EqualFold("jpg", suffix) ||
		strings.EqualFold("jpeg", suffix) || strings.EqualFold("bmp", suffix) {
		fileext = "image"
	} else if strings.EqualFold("rar", suffix) || strings.EqualFold("zip", suffix) {
		fileext = "zip"
	} else if strings.EqualFold("pdf", suffix) {
		fileext = "pdf"
	} else if strings.EqualFold("txt", suffix) {
		fileext = "text"
	} else if strings.EqualFold("xls", suffix) || strings.EqualFold("xlsx", suffix) ||
		strings.EqualFold("doc", suffix) || strings.EqualFold("docx", suffix) ||
		strings.EqualFold("ppt", suffix) || strings.EqualFold("pptx", suffix) {
		fileext = "office"
	}

	return fileext
}

// NewAttach Attach init
func NewAttach() *Attach {
	t := &Attach{}
	return t
}
