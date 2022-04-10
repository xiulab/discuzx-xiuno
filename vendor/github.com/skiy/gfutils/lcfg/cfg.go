package lcfg

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
)

var (
	name string
)

// SetCfgName SetCfgName
func SetCfgName(n string) {
	name = n
}

// GetName GetCfgName
func GetName() string {
	return name
}

// Instance Instance

// Init config init
func Init() (*gcfg.Config, error) {
	if name != "" {
		g.Config().SetFileName(name)
	}

	cfg := g.Config()
	if cfg.FilePath() == "" || cfg.GetFileName() == "" {
		return nil, fmt.Errorf("配置文件( %s )不存在", cfg.GetFileName())
	}

	return cfg, nil
}

// Get get config
func Get(n ...string) *gcfg.Config {
	return g.Config(n...)
}
