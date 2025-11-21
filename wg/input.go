package wg

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/messages"
	"math/rand"
	"time"
)

type IGraphicControl = lcl.ICustomGraphicControl

type TInput struct {
	IGraphicControl
	Text            string
	TextColor       colors.TColor
	BackgroundColor colors.TColor
	Edit            lcl.IEdit
}

func NewInput(owner lcl.IWinControl) *TInput {
	m := &TInput{IGraphicControl: lcl.NewCustomGraphicControl(owner)}
	m.Edit = lcl.NewEdit(owner)
	m.Edit.SetParent(owner)
	m.Edit.SetBorderStyle(types.BsNone)
	m.Edit.SetParentColor(true)
	m.Edit.SetLeft(-200)
	m.TextColor = colors.ClBlack
	m.BackgroundColor = colors.ClWhite
	m.SetParentBackground(true)
	m.SetParentColor(true)
	m.Canvas().SetAntialiasingMode(types.AmOn)
	m.SetControlStyle(m.ControlStyle().Include(types.CsParentBackground, types.CsFocusing))
	// 事件
	m.IGraphicControl.SetOnPaint(m.paint)
	m.IGraphicControl.SetOnWndProc(m.onWndProc)
	m.IGraphicControl.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
		m.Edit.SetFocus()
	})
	return m
}

func (m *TInput) drawBackground(canvas lcl.ICanvas) {
	canvas.BrushToBrush().SetColor(colors.RGBToColor(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))))
	canvas.PenToPen().SetColor(colors.RGBToColor(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))))
	canvas.PenToPen().SetWidth(1)
	canvas.RectangleWithIntX4(0, 0, m.Width(), m.Height())
}

func (m *TInput) drawText(canvas lcl.ICanvas) {
	clientRect := m.ClientRect()
	font := canvas.FontToFont()
	font.SetColor(m.TextColor)
	//brush := canvas.BrushToBrush()
	//brush.SetColor(m.BackgroundColor)
	canvas.TextOutWithIntX2Unicodestring(clientRect.Left, clientRect.Top, m.Text)

}

func (m *TInput) paint(sender lcl.IObject) {
	if !m.IsValid() {
		return
	}
	canvas := m.Canvas()
	m.drawBackground(canvas)
	m.drawText(canvas)
}

func (m *TInput) onWndProc(message *types.TLMessage) {
	m.InheritedWndProc(message)
	fmt.Println("msg:", message.Msg, message.Msg == messages.WM_SETFOCUS, "time:", time.Now().String())
	switch message.Msg {
	case messages.WM_SETFOCUS:

	}
}
