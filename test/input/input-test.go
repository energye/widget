package main

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"os"
	"path/filepath"
	"widget/test/util"
	"widget/wg"
)

func init() {
	util.TestLoadLibPath()
}

type TMainForm struct {
	lcl.TEngForm
	oldWndPrc uintptr
	box       lcl.IPanel
}

var MainForm TMainForm

var (
	wd, _       = os.Getwd()
	examplePath = filepath.Join(wd, "test", "input")
)

func main() {
	lcl.Init(nil, nil)
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.NewForm(&MainForm)
	lcl.Application.Run()
}

func (m *TMainForm) FormCreate(sender lcl.IObject) {
	m.SetCaption("ENERGY 自绘(自定义)输入框")
	m.SetPosition(types.PoScreenCenter)
	m.SetWidth(800)
	m.SetHeight(600)
	m.SetDoubleBuffered(true)
	m.SetColor(colors.RGBToColor(56, 57, 60))

	{
		cus := wg.NewInput(m)
		cus.SetParent(m)
		cus.SetShowHint(true)
		cus.SetCaption("上圆角")
		cus.SetHint("上圆角上圆角")
		cus.Font().SetSize(12)
		cus.Font().SetColor(colors.Cl3DFace)
	}
}
