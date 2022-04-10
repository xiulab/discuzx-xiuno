package extension

import (
	"github.com/gogf/gf/os/gcfg"
	"github.com/skiy/gfutils/lcfg"
	"github.com/skiy/gfutils/llog"
)

var (
	cfg  *gcfg.Config
	ctrl Controller
)

// Extension Extension
type Extension struct {
}

// NewExtension Extension init
func NewExtension() *Extension {
	t := &Extension{}
	return t
}

// Parsing Parsing
func (t *Extension) Parsing() {
	cfg = lcfg.Get()

	ctrl = NewGroup()
	t.ShowError(ctrl.Parsing())

	ctrl = NewUser()
	t.ShowError(ctrl.Parsing())

	ctrl = NewThreadPost()
	t.ShowError(ctrl.Parsing())

	ctrl = NewForum()
	t.ShowError(ctrl.Parsing())

	ctrl = NewFile()
	t.ShowError(ctrl.Parsing())
}

// ShowError ShowError
func (t *Extension) ShowError(err error) {
	if err != nil {
		llog.Log.Noticef("%s", err.Error())
	}
}
