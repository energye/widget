package wg

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type TInput struct {
	lcl.ICustomGraphicControl
}

func NewInput(owner lcl.IComponent) *TInput {
	m := &TInput{ICustomGraphicControl: lcl.NewCustomGraphicControl(owner)}
	m.SetWidth(120)
	m.SetHeight(40)
	m.SetParentBackground(true)
	m.SetParentColor(true)
	m.Canvas().SetAntialiasingMode(types.AmOn)
	m.SetControlStyle(m.ControlStyle().Include(types.CsParentBackground))
	return m
}
