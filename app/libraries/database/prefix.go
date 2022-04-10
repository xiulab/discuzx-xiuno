package database

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/skiy/gfutils/lcfg"
)

var (
	prefix = gmap.NewStrStrMap()
)

// InitPrefix InitPrefix
func InitPrefix() *gmap.StrStrMap {

	cfg := lcfg.Get()
	u := cfg.GetString("database.uc.0.prefix")
	d := cfg.GetString("database.discuz.0.prefix")
	x := cfg.GetString("database.xiuno.0.prefix")

	prefix.Set("uc", u)
	prefix.Set("discuz", d)
	prefix.Set("xiuno", x)

	return prefix
}

// AddPrefix AddPrefix
func AddPrefix(key string, pre string) {
	prefix.Set(key, pre)
}

// GetPrefix GetPrefix
func GetPrefix(key string) string {
	return prefix.Get(key)
}

// GetPrefixs GetPrefixs
func GetPrefixs() *gmap.StrStrMap {
	return prefix
}

// Remove Remove
func Remove(key string) {
	prefix.Remove(key)
}
