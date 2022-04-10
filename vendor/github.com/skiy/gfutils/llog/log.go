package llog

import (
	"github.com/gogf/gf/os/glog"
	"github.com/skiy/gfutils/lcfg"
)

var (
	// Log logger
	Log = glog.New()

	level = map[string]int{
		"all":           glog.LEVEL_ALL,
		"dev":           glog.LEVEL_DEV,
		"development":   glog.LEVEL_DEV,
		"prod":          glog.LEVEL_PROD,
		"production":    glog.LEVEL_PROD,
		"debug":         glog.LEVEL_DEV,
		"info":          glog.LEVEL_INFO,
		"informational": glog.LEVEL_INFO,
		"notice":        glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO,
		"warn":          glog.LEVEL_WARN,
		"warning":       glog.LEVEL_WARN,
		"error":         glog.LEVEL_ERRO,
		"critical":      glog.LEVEL_CRIT,
		"alert":         glog.LEVEL_PROD | glog.LEVEL_NOTI | glog.LEVEL_INFO,
	}
)

// Instance Instance
func Instance() *glog.Logger {
	log := glog.New()
	log.SetLevel(glog.LEVEL_ALL)
	return log
}

// InitLog log init
func InitLog() error {
	cfg := lcfg.Get()

	//日志等级
	if logLevel := cfg.GetString("log.level"); logLevel != "" {
		if l, ok := level[logLevel]; ok {
			Log.SetLevel(l)
		}
	}

	//日志路径
	if logPath := cfg.GetString("log.path"); logPath != "" {
		if err := Log.SetPath(logPath); err != nil {
			Log.Warning(err.Error())
		}
	}

	// 是否输出错误
	Log.Stack(cfg.GetBool("log.trace"))

	return nil
}

// GetLog GetLog
func GetLog() *glog.Logger {
	return Log
}
