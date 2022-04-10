package app

import (
	"discuzx-xiuno/app/controllers"
	"discuzx-xiuno/app/controllers/extension"
	"errors"
	"fmt"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
)

// App App
type App struct {
}

// NewApp App init
func NewApp() *App {
	t := &App{}
	return t
}

// Parsing Parsing
func (t *App) Parsing() {
	tablesName := [...]string{
		"user",
		"group",
		"forum",
		"attach",
		"thread",
		"post",
		"thread_top",
		"mythread",
		"mypost",
	}

	var err error
	var ctrl controllers.Controller

	// 遍历表
	for _, table := range tablesName {
		cfgOffset := fmt.Sprintf("tables.xiuno.%s", table)

		// 是否转换
		if lcfg.Get().GetBool(fmt.Sprintf("%s.convert", cfgOffset)) {
			// 转换控制器
			if ctrl, err = t.ctrl(table); err != nil {
				llog.Log.Noticef("%s(%s)", err.Error(), table)
				continue
			}

			if err = ctrl.ToConvert(); err != nil {
				llog.Log.Fatalf("转换数据表(%s)失败: %s", table, err.Error())
			}
		}
	}

	// 扩展功能
	extension.NewExtension().Parsing()
}

// ctrl 返回控制器
func (t *App) ctrl(name string) (ctrl controllers.Controller, err error) {
	switch name {
	case "user":
		ctrl = controllers.NewUser()

	case "group":
		ctrl = controllers.NewGroup()

	case "forum":
		ctrl = controllers.NewForum()

	case "attach":
		ctrl = controllers.NewAttach()

	case "thread":
		ctrl = controllers.NewThread()

	case "post":
		ctrl = controllers.NewPost()

	case "thread_top":
		ctrl = controllers.NewThreadTop()

	case "mythread":
		ctrl = controllers.NewMythread()

	case "mypost":
		ctrl = controllers.NewMypost()

	default:
		err = errors.New("找不到对应的控制器,无法转换该表")

	}

	return
}
