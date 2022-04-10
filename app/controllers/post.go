package controllers

import (
	"database/sql"
	"discuzx-xiuno/app/libraries/common"
	"discuzx-xiuno/app/libraries/database"
	"fmt"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
	"strings"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/util/gconv"
	"github.com/skiy/bbcode"
)

// Post Post
type Post struct {
}

// ToConvert ToConvert
func (t *Post) ToConvert() (err error) {
	start := time.Now()

	cfg := lcfg.Get()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxPostTable := discuzPre + "forum_post"

	lastPid := cfg.GetInt("tables.xiuno.post.last_pid")

	fields := "tid,pid,authorid,first,dateline,useip,message,attachment"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxPostTable+" t").Where("pid >= ?", lastPid).OrderBy("pid ASC").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
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

	var count int64
	batch := cfg.GetInt("tables.xiuno.post.batch")

	dataList := gdb.List{}
	countMax := len(r.List())

	for _, u := range r.List() {
		userip := common.IP2Long(gconv.String(u["useip"]))
		messageFmt := gconv.String(u["message"])

		if messageFmt != "" {
			messageFmt = t.BBCodeToHTML(messageFmt) //处理message中的附件
		}

		// 添加图片附件到 messageFmt
		if gconv.Int64(u["attachment"]) > 0 {
			messageFmt += t.addImageToContent(gconv.Int64(u["pid"]))
		}

		d := gdb.Map{
			"tid":         u["tid"],
			"pid":         u["pid"],
			"uid":         u["authorid"],
			"isfirst":     u["first"],
			"create_date": u["dateline"],
			"userip":      userip,
			"message":     messageFmt,
			"message_fmt": messageFmt,
		}

		// 批量插入数量
		if batch > 1 {
			dataList = append(dataList, d)
		} else {
			res, err := xiunoDB.Insert(xiunoTable, d)
			if err != nil {
				return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
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

	llog.Log.Infof("表 %s 数据导入成功, 本次导入: %d/%d 条数据, 耗时: %v", xiunoTable, count, countMax, time.Since(start))
	return
}

// NewPost Post init
func NewPost() *Post {
	t := &Post{}
	return t
}

// BBCodeToHTML bbcode 转 html
func (t *Post) BBCodeToHTML(msg string) string {
	compiler := bbcode.NewCompiler(true, true)

	//转 table
	compiler.SetTag("table", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "table"
		return out, true
	})

	//转 tr
	compiler.SetTag("tr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "tr"

		return out, true
	})
	//转 td
	compiler.SetTag("td", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "td"

		return out, true
	})

	//ul
	compiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "ul"

		return out, true
	})

	//text-align
	compiler.SetTag("align", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "div"
		value := node.GetOpeningTag().Value
		if value != "" {
			out.Attrs["style"] = "text-align: " + value
		}
		return out, true
	})

	//backcolor=yellow
	compiler.SetTag("backcolor", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "span"
		value := node.GetOpeningTag().Value
		if value != "" {
			out.Attrs["style"] = "background-color: " + value
		}
		return out, true
	})

	//li -> 将 [*] 转为 li
	compiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "li"

		return out, true
	})

	//font -> 清空 font
	compiler.SetTag("font", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		return out, true
	})

	//free -> 清空 free
	compiler.SetTag("free", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		return out, true
	})

	//hide -> 清空 hide
	compiler.SetTag("hide", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		return out, true
	})

	//qq -> 更新 qq 标签
	compiler.SetTag("qq", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		value := node.GetOpeningTag().Value
		if value == "" {
			qq := bbcode.CompileText(node)

			if len(qq) > 0 {
				out.Name = "a"
				out.Attrs["href"] = "http://wpa.qq.com/msgrd?v=3&uin=" + qq + "&site=Xiuno&from=Xiuno&menu=yes"
				out.Attrs["target"] = "_blank"
			}
		}
		return out, true
	})

	//处理message中的附件
	xiunoTable := database.GetPrefix("xiuno") + lcfg.Get().GetString("tables.xiuno.attach.name")

	compiler.SetTag("attach", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {

		out := bbcode.NewHTMLTag("")
		out.Name = ""

		closeFlag := true

		value := node.GetOpeningTag().Value
		if value == "" {
			attachID := bbcode.CompileText(node)

			if len(attachID) > 0 {
				r, err := database.GetXiunoDB().Table(xiunoTable).Where("aid = ?", attachID).Fields("isimage,filename").One()
				if err != nil {
					if err != sql.ErrNoRows {
						llog.Log.Noticef("查询附件(aid: %s)失败, %s", attachID, err.Error())
					}
				} else if r != nil {

					isimage := r["isimage"].Int()
					if isimage == 1 {
						out.Name = "img"
						out.Attrs["src"] = "upload/attach/" + r["filename"].String()

						closeFlag = false
					} else {
						out.Name = "a"
						out.Attrs["href"] = "?attach-download-" + attachID + ".htm" //bbcode.ValidURL(filename)
						out.Attrs["target"] = "_blank"
						out.Value = "附件: "

						closeFlag = true
					}
				}
			}
		}
		return out, closeFlag
	})

	return compiler.Compile(msg)
}

// addImageToContent 添加图片到内容里
func (t *Post) addImageToContent(pid int64) string {
	xiunoTable := database.GetPrefix("xiuno") + lcfg.Get().GetString("tables.xiuno.attach.name")
	r, err := database.GetXiunoDB().Table(xiunoTable).Where("pid = ? AND isimage = 1", pid).OrderBy("create_date ASC").All()
	if err != nil {
		if err == sql.ErrNoRows {
			llog.Log.Debugf("表 %s 无帖子(%d)的图片附件", xiunoTable, pid)
			return ""
		}
		llog.Log.Debugf("表 %s 查询帖子(%d)图片附件失败, %s", xiunoTable, pid, err.Error())
		return ""
	}

	var imagesStrArr []string
	var oriWidth, oriHeigth, proportion float64
	var style string
	for _, u := range r.List() {
		imageURL := "upload/attach/" + gconv.String(u["filename"])

		oriWidth = gconv.Float64(u["width"])
		oriHeigth = gconv.Float64(u["height"])
		if oriWidth == 0 {
			oriWidth = 400
		} else if oriWidth > 800 {
			proportion = 800 / oriWidth
			oriWidth = 800
			oriHeigth = oriHeigth * proportion
			style = " style=\"width: 800px; height: auto; cursor: pointer;\""
		}

		imageStr := fmt.Sprintf("<a href=\"%s\" target=\"_blank\" _href=\"%s\"><img src=\"%s\" width=\"%d\" height=\"%d\" %s></a>",
			imageURL, imageURL, imageURL, int64(oriWidth), int64(oriHeigth), style)
		imagesStrArr = append(imagesStrArr, imageStr)
	}

	return "<p class=\"imagelist\">" + strings.Join(imagesStrArr, "<br>") + "</p>"
}
