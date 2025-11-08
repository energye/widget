package wg

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"path/filepath"
	"strings"
)

type RoundedCorner = int32

const (
	RcLeftTop RoundedCorner = iota
	RcRightTop
	RcLeftBottom
	RcRightBottom
)

type TRoundedCorners = types.TSet

const iconMargin = 5

type buttonState = int32

const (
	bsDefault  buttonState = iota // 默认状态
	bsEnter                       // 移入状态
	bsDown                        // 按下状态
	bsDisabled                    // 禁用状态
)

var (
	defaultButtonColor        = colors.RGBToColor(66, 133, 244)  // 淡蓝色
	defaultButtonColorDisable = colors.RGBToColor(200, 200, 200) // 浅灰色
)

// TButton 多功能自绘按钮
// 颜色状态: 默认颜色, 移入颜色, 按下颜色, 禁用颜色
// 当大小改变, 颜色改变 会重新绘制
type TButton struct {
	lcl.ICustomGraphicControl
	isDisable                bool            // 是否禁用
	alpha                    byte            // 透明度 0 ~ 255
	radius                   int32           // 圆角度
	autoSize                 bool            // 自动大小
	text                     string          // 文本
	RoundedCorner            TRoundedCorners // 按钮圆角方向，默认四角
	TextOffSetX, TextOffSetY int32           // 文本显示偏移位置
	// 图标
	iconFavorite       lcl.IPicture // 按钮前置图标, 靠左
	iconClose          lcl.IPicture // 按钮关闭图标, 靠右
	iconCloseHighlight lcl.IPicture // 按钮关闭图标移入高亮, 靠右
	isEnterClose       bool         // 鼠标是否移入关闭图标
	icon               lcl.IPicture // 按钮图标, 中间
	// 用户事件
	onCloseClick lcl.TNotifyEvent
	onPaint      lcl.TNotifyEvent
	onMouseEnter lcl.TNotifyEvent
	onMouseLeave lcl.TNotifyEvent
	onClick      lcl.TNotifyEvent
	onMouseDown  lcl.TMouseEvent
	onMouseUp    lcl.TMouseEvent
	// 默认颜色, 移入颜色, 按下颜色, 禁用颜色
	buttonState   buttonState
	defaultColor  *TButtonColor
	enterColor    *TButtonColor
	downColor     *TButtonColor
	disabledColor *TButtonColor
}

func NewButton(owner lcl.IComponent) *TButton {
	m := &TButton{ICustomGraphicControl: lcl.NewCustomGraphicControl(owner)}
	m.SetWidth(120)
	m.SetHeight(40)
	m.SetParentBackground(true)
	m.SetParentColor(true)
	m.Canvas().SetAntialiasingMode(types.AmOn)
	m.SetControlStyle(m.ControlStyle().Include(types.CsParentBackground))
	m.alpha = 255
	m.radius = 0
	m.ICustomGraphicControl.SetOnPaint(m.paint)
	m.ICustomGraphicControl.SetOnMouseEnter(m.enter) // 进入
	m.ICustomGraphicControl.SetOnMouseLeave(m.leave) // 移出
	m.ICustomGraphicControl.SetOnMouseDown(m.down)   // 按下
	m.ICustomGraphicControl.SetOnMouseUp(m.up)       // 抬起
	m.ICustomGraphicControl.SetOnMouseMove(m.move)
	m.RoundedCorner = types.NewSet(RcLeftTop, RcRightTop, RcLeftBottom, RcRightBottom)
	m.iconFavorite = lcl.NewPicture()
	m.iconClose = lcl.NewPicture()
	m.iconCloseHighlight = lcl.NewPicture()
	m.icon = lcl.NewPicture()
	m.iconFavorite.SetOnChange(m.iconChange)
	m.iconClose.SetOnChange(m.iconChange)
	m.iconCloseHighlight.SetOnChange(m.iconChange)
	m.icon.SetOnChange(m.iconChange)
	// 创建按钮颜色对象
	m.defaultColor = NewButtonColor()
	m.defaultColor.type_ = bsDefault
	m.enterColor = NewButtonColor()
	m.enterColor.type_ = bsEnter
	m.downColor = NewButtonColor()
	m.downColor.type_ = bsDown
	m.disabledColor = NewButtonColor()
	m.disabledColor.type_ = bsDisabled
	m.SetDisabledColor(defaultButtonColorDisable, defaultButtonColorDisable)

	m.SetColor(defaultButtonColor)
	m.SetBorderColor(defaultButtonColor)
	m.SetBorderDirections(types.NewSet(BbdLeft, BbdTop, BbdRight, BbdBottom))

	// 销毁事件
	m.SetOnDestroy(func() {
		//fmt.Println("Graphic Button 释放资源")
		// 清空事件
		m.ICustomGraphicControl.SetOnPaint(nil)
		m.ICustomGraphicControl.SetOnMouseEnter(nil)
		m.ICustomGraphicControl.SetOnMouseLeave(nil)
		m.ICustomGraphicControl.SetOnMouseDown(nil)
		m.ICustomGraphicControl.SetOnMouseUp(nil)
		m.ICustomGraphicControl.SetOnMouseMove(nil)
		m.iconFavorite.SetOnChange(nil)
		m.iconClose.SetOnChange(nil)
		m.iconCloseHighlight.SetOnChange(nil)
		m.icon.SetOnChange(nil)
		m.SetOnDestroy(nil)
		// 释放持有资源
		m.iconFavorite.Free()
		m.iconClose.Free()
		m.iconCloseHighlight.Free()
		m.icon.Free()
		//m.imgPool.Free()
		//m.imgBitmapPool.Free()
	})
	return m
}

func (m *TButton) enter(sender lcl.IObject) {
	if m.isDisable || !m.IsValid() {
		return
	}
	m.buttonState = bsEnter
	m.Invalidate()
	if m.onMouseEnter != nil {
		m.onMouseEnter(sender)
	}
}

func (m *TButton) leave(sender lcl.IObject) {
	if m.isDisable || !m.IsValid() {
		return
	}
	m.isEnterClose = false
	m.buttonState = bsDefault
	m.Invalidate()
	if m.onMouseLeave != nil {
		m.onMouseLeave(sender)
	}
}

func (m *TButton) down(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.isDisable || !m.IsValid() {
		return
	}
	if m.isCloseArea(X, Y) {
		if m.onCloseClick != nil {
			m.onCloseClick(sender)
		}
	} else {
		m.buttonState = bsDown
		m.Invalidate()
		if m.onMouseDown != nil {
			m.onMouseDown(sender, button, shift, X, Y)
		}
	}
}

func (m *TButton) up(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.isDisable || !m.IsValid() {
		return
	}
	m.buttonState = bsEnter
	m.Invalidate()
	if m.onMouseUp != nil {
		m.onMouseUp(sender, button, shift, X, Y)
	}
}

func (m *TButton) SetDisable(disable bool) {
	m.isDisable = disable
	if m.isDisable {
		m.buttonState = bsDisabled
	} else {
		m.buttonState = bsDefault
	}
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.Invalidate()
	})
}
func (m *TButton) iconChange(sender lcl.IObject) {
	if m.isDisable || !m.IsValid() {
		return
	}
	m.Invalidate()
}

func (m *TButton) isCloseArea(X int32, Y int32) bool {
	if m.isDisable || !m.IsValid() {
		return false
	}
	rect := m.ClientRect()
	closeX := rect.Width() - m.iconClose.Width() - iconMargin
	closeY := rect.Height()/2 - m.iconClose.Height()/2
	return X >= closeX && X <= rect.Width()-iconMargin && Y >= closeY && Y <= rect.Height()/2+m.iconClose.Height()
}

func (m *TButton) move(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
	if m.isDisable || !m.IsValid() {
		return
	}
	lcl.Screen.SetCursor(types.CrDefault)
	if m.isCloseArea(X, Y) {
		if !m.isEnterClose {
			m.isEnterClose = true
			m.Invalidate()
		}
	} else if m.isEnterClose {
		m.isEnterClose = false
		m.Invalidate()
	}
}

func (m *TButton) drawRoundedGradientButton(canvas lcl.ICanvas, rect types.TRect) {
	text := m.text
	var color *TButtonColor
	switch m.buttonState {
	case bsDefault:
		color = m.defaultColor
	case bsEnter:
		color = m.enterColor
	case bsDown:
		color = m.downColor
	case bsDisabled:
		color = m.disabledColor
	}
	if color == nil {
		return
	}
	color.paint(m.RoundedCorner, rect, m.alpha, m.radius)

	// 绘制到目标画布
	canvas.DrawWithIntX2Graphic(rect.Left, rect.Top, color.bitMap)

	// 绘制按钮文字（在原始画布上绘制，确保文字不透明）
	brush := canvas.BrushToBrush()
	brush.SetStyle(types.BsClear)

	textMargin := int32(0) // 文本与图标的间距
	// 计算左右图标占用的空间
	leftArea := int32(0)
	if m.iconFavorite.Width() > 0 {
		leftArea = iconMargin + m.iconFavorite.Width() + iconMargin // 左边距10 + 图标宽度 + 图标与文本间距10
		textMargin += iconMargin
	}

	rightArea := int32(0)
	if m.iconClose.Width() > 0 {
		rightArea = iconMargin + m.iconClose.Width() + iconMargin // 右边距10 + 图标宽度 + 图标与文本间距10
		textMargin += -iconMargin
	}

	// 计算文本可用宽度
	availWidth := rect.Width() - leftArea - rightArea
	if availWidth < 0 {
		availWidth = 0
	}
	// 截断文本
	if len(text) > 0 {
		text = truncateText(canvas, text, availWidth)
	}

	// 计算文字位置
	textSize := canvas.TextExtentWithUnicodestring(text)
	textX := rect.Left + m.TextOffSetX + (rect.Width()-textSize.Cx)/2 + textMargin
	textY := rect.Top + m.TextOffSetY + (rect.Height()-textSize.Cy)/2

	// 绘制文字阴影（增强可读性）
	//canvas.FontToFont().SetColor(colors.ClBlack)
	//canvas.TextOutWithIntX2Unicodestring(textX+1, textY+1, text)

	// 绘制主文字
	//canvas.FontToFont().SetColor(colors.ClWhite)
	canvas.TextOutWithIntX2Unicodestring(textX, textY, text)

	// 左: 绘制图标 favorite
	favY := rect.Height()/2 - m.iconFavorite.Height()/2
	canvas.DrawWithIntX2Graphic(iconMargin, favY, m.iconFavorite.Graphic())

	// 右: 绘制图标 close
	iconClose := m.iconClose
	if m.isEnterClose {
		iconClose = m.iconCloseHighlight
	}
	closeX := rect.Width() - iconClose.Width() - iconMargin
	closeY := rect.Height()/2 - iconClose.Height()/2
	canvas.DrawWithIntX2Graphic(closeX, closeY, iconClose.Graphic())

	// 中间: 绘制图标 icon
	iconW, iconH := m.icon.Width(), m.icon.Height()
	iconX := rect.Left + (rect.Width()-iconW)/2
	iconY := rect.Top + (rect.Height()-iconH)/2
	canvas.DrawWithIntX2Graphic(iconX, iconY, m.icon.Graphic())
}

func (m *TButton) Disable() bool {
	return m.isDisable
}

func (m *TButton) SetCaption(value string) {
	m.SetText(value)
}

func (m *TButton) Caption() string {
	return m.text
}

func (m *TButton) SetText(value string) {
	m.text = value
	if m.autoSize {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			// 自动大小, 根据文本宽自动调整按钮宽度
			leftArea := int32(0)
			if m.iconFavorite.Width() > 0 {
				leftArea = iconMargin + m.iconFavorite.Width() + iconMargin
			}
			rightArea := int32(0)
			if m.iconClose.Width() > 0 {
				rightArea = iconMargin + m.iconClose.Width() + iconMargin
			}
			textWidth := m.Canvas().TextWidthWithUnicodestring(m.text)
			width := textWidth + leftArea + rightArea + iconMargin*2
			m.SetWidth(width)
		})
	} else {
		m.Invalidate()
	}
}

func (m *TButton) Text() string {
	return m.text
}

func (m *TButton) SetAutoSize(v bool) {
	m.autoSize = v
}

func (m *TButton) SetIconFavorite(filePath string) {
	if !m.IsValid() {
		return
	}
	m.iconFavorite.LoadFromFile(filePath)
}

func (m *TButton) SetIconFavoriteFormBytes(pngData []byte) {
	if !m.IsValid() {
		return
	}
	mem := lcl.NewMemoryStream()
	defer mem.Free()
	lcl.StreamHelper.WriteBuffer(mem, pngData)
	mem.SetPosition(0)
	m.iconFavorite.LoadFromStream(mem)
}

func (m *TButton) SetIcon(filePath string) {
	if !m.IsValid() {
		return
	}
	m.icon.LoadFromFile(filePath)
}

func (m *TButton) SetIconFormBytes(pngData []byte) {
	if !m.IsValid() {
		return
	}
	mem := lcl.NewMemoryStream()
	defer mem.Free()
	lcl.StreamHelper.WriteBuffer(mem, pngData)
	mem.SetPosition(0)
	m.icon.LoadFromStream(mem)
}

func (m *TButton) SetIconClose(filePath string) {
	if !m.IsValid() {
		return
	}
	path, name := filepath.Split(filePath)
	ns := strings.Split(name, ".")
	enterFilePath := filepath.Join(path, ns[0]+"_enter.png")
	m.iconClose.LoadFromFile(filePath)
	m.SetIconCloseHighlight(enterFilePath)
}

func (m *TButton) SetIconCloseHighlight(filePath string) {
	if !m.IsValid() {
		return
	}
	m.iconCloseHighlight.LoadFromFile(filePath)
}

func (m *TButton) SetIconCloseFormBytes(pngData []byte) {
	if !m.IsValid() {
		return
	}
	mem := lcl.NewMemoryStream()
	defer mem.Free()
	lcl.StreamHelper.WriteBuffer(mem, pngData)
	mem.SetPosition(0)
	m.iconClose.LoadFromStream(mem)
}

func (m *TButton) SetIconCloseHighlightFormBytes(pngData []byte) {
	if !m.IsValid() {
		return
	}
	mem := lcl.NewMemoryStream()
	defer mem.Free()
	lcl.StreamHelper.WriteBuffer(mem, pngData)
	mem.SetPosition(0)
	m.iconCloseHighlight.LoadFromStream(mem)
}

func (m *TButton) paint(sender lcl.IObject) {
	if !m.IsValid() {
		return
	}
	m.drawRoundedGradientButton(m.Canvas(), m.ClientRect())
	if m.onPaint != nil {
		m.onPaint(sender)
	}
}
func (m *TButton) SetOnCloseClick(fn lcl.TNotifyEvent) {
	m.onCloseClick = fn
}

func (m *TButton) SetOnPaint(fn lcl.TNotifyEvent) {
	m.onPaint = fn
}

func (m *TButton) SetOnMouseDown(fn lcl.TMouseEvent) {
	m.onMouseDown = fn
}

func (m *TButton) SetOnMouseUp(fn lcl.TMouseEvent) {
	m.onMouseUp = fn
}

func (m *TButton) SetOnMouseEnter(fn lcl.TNotifyEvent) {
	m.onMouseEnter = fn
}

func (m *TButton) SetOnMouseLeave(fn lcl.TNotifyEvent) {
	m.onMouseLeave = fn
}

// SetDefaultColor 设置按钮的默认颜色
// start: 按钮默认状态下的起始颜色
// end: 按钮默认状态下的结束颜色
func (m *TButton) SetDefaultColor(start, end colors.TColor) {
	// 更新按钮默认颜色配置
	m.defaultColor.start = start
	m.defaultColor.end = end
	m.defaultColor.forcePaint(m.RoundedCorner, m.ClientRect(), m.alpha, m.radius)
}

func (m *TButton) DefaultColor() (start, end colors.TColor) {
	start = m.defaultColor.start
	end = m.defaultColor.end
	return
}

func (m *TButton) SetEnterColor(start, end colors.TColor) {
	m.enterColor.start = start
	m.enterColor.end = end
	m.enterColor.forcePaint(m.RoundedCorner, m.ClientRect(), m.alpha, m.radius)
}

func (m *TButton) EnterColor() (start, end colors.TColor) {
	start = m.enterColor.start
	end = m.enterColor.end
	return
}

func (m *TButton) SetDownColor(start, end colors.TColor) {
	m.downColor.start = start
	m.downColor.end = end
	m.downColor.forcePaint(m.RoundedCorner, m.ClientRect(), m.alpha, m.radius)
}

func (m *TButton) DownColor() (start, end colors.TColor) {
	start = m.downColor.start
	end = m.downColor.end
	return
}

func (m *TButton) SetDisabledColor(start, end colors.TColor) {
	m.disabledColor.start = start
	m.disabledColor.end = end
	m.disabledColor.forcePaint(m.RoundedCorner, m.ClientRect(), m.alpha, m.radius)
}

func (m *TButton) DisabledColor() (start, end colors.TColor) {
	start = m.disabledColor.start
	end = m.disabledColor.end
	return
}

func (m *TButton) SetColor(color colors.TColor) {
	m.SetColorGradient(color, color)
}

// SetColorGradient 设置按钮的颜色渐变效果
// start: 渐变起始颜色
// end: 渐变结束颜色
func (m *TButton) SetColorGradient(start, end colors.TColor) {
	m.SetDefaultColor(start, end)
	m.SetEnterColor(DarkenColor(start, 0.1), DarkenColor(end, 0.1))
	m.SetDownColor(DarkenColor(start, 0.2), DarkenColor(end, 0.2))
}

// SetBorderColor 设置按钮所有状态下的边框颜色
//
// 参数:
//
//	color - 指定的边框颜色值
//
// 该函数会同时设置按钮在默认、悬停、按下和禁用状态下的边框颜色
// 为统一的颜色值，实现按钮边框颜色的整体变更
func (m *TButton) SetBorderColor(color colors.TColor) {
	m.defaultColor.Border.Color = color
	m.enterColor.Border.Color = DarkenColor(color, 0.1)
	m.downColor.Border.Color = DarkenColor(color, 0.2)
}

// SetBorderDirections 设置按钮的所有状态边框样式
// 参数:
//
//	borders: TButtonBorders类型，指定按钮的边框样式
//
// 该函数会同时设置按钮的默认、悬停、按下和禁用四种状态的边框样式
// 为相同的值，实现统一的边框外观效果
func (m *TButton) SetBorderDirections(directions TButtonBorderDirections) {
	m.defaultColor.Border.Direction = directions
	m.enterColor.Border.Direction = directions
	m.downColor.Border.Direction = directions
	m.disabledColor.Border.Direction = directions
}

func (m *TButton) SetAlpha(alpha byte) {
	m.alpha = alpha
}

func (m *TButton) SetRadius(radius int32) {
	m.radius = radius
}

func (m *TButton) Free() {
	m.ICustomGraphicControl.Free()
}

// 文本截断函数（添加在文本末尾）
func truncateText(canvas lcl.ICanvas, text string, maxWidth int32) string {
	if maxWidth <= 0 {
		return ""
	}
	ellipsis := "..."
	ellipsisWidth := canvas.TextWidthWithUnicodestring(ellipsis)
	if ellipsisWidth > maxWidth {
		return ""
	}
	textWidth := canvas.TextWidthWithUnicodestring(text)
	if textWidth <= maxWidth {
		return text
	}
	// 二分查找截断位置
	runes := []rune(text)
	left, right := 0, len(runes)
	for left < right {
		mid := (left + right) / 2
		truncated := string(runes[:mid]) + ellipsis
		if canvas.TextWidthWithUnicodestring(truncated) <= maxWidth {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left == 0 {
		return ellipsis
	}
	return string(runes[:left-1]) + ellipsis
}
