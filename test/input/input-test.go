package main

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/test/util"
	"github.com/energye/widget/wg"
	"os"
	"path/filepath"
	"time"
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
		var inputs []*wg.TInput
		newInput := func(h int32) {
			for i := int32(0); i < 10; i++ {
				cus := wg.NewInput(m)
				cus.Font().SetSize(12)
				cus.Font().SetColor(colors.Cl3DFace)
				cus.SetWidth(150)
				cus.SetHeight(40)
				cus.Text = time.Now().String()
				cus.SetLeft(i*cus.Width() + 10)
				cus.SetTop(h + 10)
				cus.SetParent(m)
				inputs = append(inputs, cus)
			}
		}
		for h := int32(0); h < 15; h++ {
			newInput(h * 40)
		}
		go func() {
			for {
				time.Sleep(time.Second / 2)
				for _, input := range inputs {
					go input.Invalidate()
				}
			}
		}()
	}
}
