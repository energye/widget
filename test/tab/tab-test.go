package main

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	examplePath = filepath.Join(wd, "test", "tab")
)

func main() {
	lcl.Init(nil, nil)
	lcl.Application.Initialize()
	lcl.Application.SetMainFormOnTaskBar(true)
	lcl.Application.NewForm(&MainForm)
	lcl.Application.Run()
}

func (m *TMainForm) FormCreate(sender lcl.IObject) {
	m.SetCaption("ENERGY 自绘(自定义) Tab ")
	m.SetPosition(types.PoScreenCenter)
	m.SetWidth(800)
	m.SetHeight(600)
	m.SetDoubleBuffered(true)
	//m.SetColor(colors.RGBToColor(56, 57, 60))

	box := lcl.NewPanel(m)
	box.SetTop(100)
	box.SetWidth(800)
	box.SetHeight(500)
	box.SetBevelInner(types.BvNone)
	box.SetBevelOuter(types.BvNone)
	box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
	box.SetParent(m)

	tab := wg.NewTab(box)
	tab.SetParent(box)
	tab.SetBounds(0, 0, box.Width(), box.Height())
	tab.SetAlign(types.AlClient)
	//tab.EnableScrollButton(false)
	tab.RecalculatePosition()

	addPage := func(count int) {
		page := tab.NewPage()
		page.SetColor(colors.RGBToColor(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))))
		testPanel := lcl.NewPanel(page)
		testPanel.SetColor(colors.RGBToColor(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))))
		testPanel.SetTop(int32(rand.Intn(400)))
		testPanel.SetLeft(int32(rand.Intn(400)))
		testPanel.SetParent(page)
		btn := page.Button()
		btn.SetText(RandMixString())
		btn.SetIconFavorite("C:\\app\\workspace\\widget\\test\\tab\\resources\\icon.png")
		btn.SetIconClose("C:\\app\\workspace\\widget\\test\\tab\\resources\\close.png")
		testButton := wg.NewButton(page)
		testButton.SetLeft(20)
		testButton.SetTop(20)
		testButton.SetRadius(8)
		testButton.SetAutoSize(true)
		testButton.SetText("自绘按钮" + RandMixString())
		testButton.Font().SetColor(colors.ClWhite)
		testButton.Font().SetStyle(types.NewSet(types.FsBold))
		testButton.SetParent(page)
	}

	count := 0
	add := wg.NewButton(m)
	add.SetLeft(20)
	add.SetTop(20)
	add.SetRadius(8)
	add.SetText("添加一个 Tab Page")
	add.Font().SetColor(colors.ClWhite)
	add.Font().SetStyle(types.NewSet(types.FsBold))
	add.SetOnClick(func(sender lcl.IObject) {
		addPage(count)
		count++
	})
	add.SetParent(m)

	for i := 0; i < 10; i++ {
		addPage(count)
		count++
	}
	lcl.RunOnMainThreadAsync(func(id uint32) {
		tab.RecalculatePosition()
	})
}

func RandMixString() string {
	rand.Seed(time.Now().UnixNano())
	hanzi := []rune("的一是在不了有和人这中大为上个国我以要他时来用们")
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	length := rand.Intn(10)
	var builder strings.Builder
	builder.WriteRune(hanzi[rand.Intn(len(hanzi))])
	builder.WriteByte(letters[rand.Intn(len(letters))])
	for i := 2; i < length; i++ {
		if rand.Intn(2) == 0 {
			builder.WriteRune(hanzi[rand.Intn(len(hanzi))])
		} else {
			builder.WriteByte(letters[rand.Intn(len(letters))])
		}
	}
	return builder.String()
}
