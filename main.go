package main

import (
	"discuzx-xiuno/app"
	"discuzx-xiuno/app/libraries/database"
	"flag"
	"fmt"
	"github.com/gogf/gf/os/gcfg"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
	"time"
)

var (
	err error
	cfg *gcfg.Config
)

const (
	version = "2.0.3"
	verDate = "2020/02/01"
)

// 配置初始化
func init() {
	cfgName := loadCfgName()
	if cfgName != "" {
		lcfg.SetCfgName(cfgName)
	}

	cfg, err = lcfg.Init()
	if err != nil {
		return
	}

	err = llog.InitLog()
	if err != nil {
		return
	}

	//初始化前缀
	database.InitPrefix()
}

// loadCfgName 加载配置文件名
// 运行时加参数 -config=config.dev.toml 指定配置文件
func loadCfgName() string {
	configFile := flag.String("config", "config.toml", "-config=config.dev.toml")
	flag.Parse()
	return *configFile
}

func main() {
	fmt.Printf(`
:::
::: Discuz!X 3.x 转 XiunoBBS 4.x 工具
:::
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
::: 本项目讨论帖：https://bbs.jadehive.com/thread-8059.htm
:::
::: 执行过程中按 "Ctrl + Z" 结束本程序...
:::
::: 版本: %s    时间: %s
:::

`, version, verDate)

	// 判断 MYSQL 连接是否正常
	if err := checkConnectDB(); err != nil {
		llog.Log.Fatalf("数据库连接失败: %s", err.Error())
	}

	start := time.Now()

	llog.Log.Info("开始导入数据 ...")

	app.NewApp().Parsing()

	llog.Log.Infof("已将 Discuz!X 转换至 XiunoBBS, 总耗时: %v\n", time.Since(start))

	fmt.Printf(`
:::
::: 本项目开源地址: https://github.com/skiy/discuzx-xiuno
::: 技术支持论坛: https://bbs.jadehive.com
::: 如需技术支持请加 QQ 群: 891844359
:::
::: 如有意见和建议或者遇到 BUG，请到 GitHub 提 issue 。
:::

`)

}

// checkConnectDB 检测数据库连接是否正常
func checkConnectDB() (err error) {
	if err = database.GetDiscuzDB().PingMaster(); err != nil {
		return fmt.Errorf("%s(Discuz!X)", err.Error())
	}

	if err = database.GetUcDB().PingMaster(); err != nil {
		return fmt.Errorf("%s(Discuz!UCenter)", err.Error())
	}

	if err = database.GetXiunoDB().PingMaster(); err != nil {
		return fmt.Errorf("%s(XiunoBBS)", err.Error())
	}

	return
}
