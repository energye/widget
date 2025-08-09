package wg

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"math"
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

type RoundedCorners = types.TSet

type TButton struct {
	lcl.ICustomGraphicControl
	startColor               colors.TColor  // 按钮起始渐变颜色
	endColor                 colors.TColor  // 按钮结束渐变颜色
	activeColor              float32        // 按钮激活颜色深度 0.0 ~ 1.0
	alpha                    byte           // 透明度 0 ~ 255
	radius                   int32          // 圆角度
	isEnter                  bool           // 鼠标是否移入
	isDown                   bool           // 鼠标是否按下
	RoundedCorner            RoundedCorners // 按钮圆角方向，默认四角
	TextOffSetX, TextOffSetY int32          // 文本显示偏移位置
	// 图标
	iconFavorite       lcl.IPicture // 按钮前置图标
	iconClose          lcl.IPicture // 按钮关闭图标
	iconCloseHighlight lcl.IPicture // 按钮关闭图标 高亮
	isEnterClose       bool         // 鼠标是否移入关闭图标
	icon               lcl.IPicture // 按钮图标
	// 用户事件
	onCloseClick lcl.TNotifyEvent
	onPaint      lcl.TNotifyEvent
	onMouseEnter lcl.TNotifyEvent
	onMouseLeave lcl.TNotifyEvent
	onMouseDown  lcl.TMouseEvent
	onMouseUp    lcl.TMouseEvent
	// 是否禁用
	IsDisable bool
}

func NewButton(owner lcl.IComponent) *TButton {
	m := &TButton{ICustomGraphicControl: lcl.NewCustomGraphicControl(owner)}
	m.SetWidth(120)
	m.SetHeight(40)
	m.SetParentBackground(true)
	m.SetParentColor(true)
	m.Canvas().SetAntialiasingMode(types.AmOn)
	m.SetControlStyle(m.ControlStyle().Include(types.CsParentBackground))
	m.startColor = colors.ClBlue
	m.endColor = colors.ClNavy
	m.alpha = 180
	m.radius = 10
	m.ICustomGraphicControl.SetOnPaint(m.paint)
	m.ICustomGraphicControl.SetOnMouseEnter(m.enter)
	m.ICustomGraphicControl.SetOnMouseLeave(m.leave)
	m.ICustomGraphicControl.SetOnMouseDown(m.down)
	m.ICustomGraphicControl.SetOnMouseUp(m.up)
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
	m.SetOnDestroy(func() {
		fmt.Println("释放资源")
		m.iconFavorite.Free()
		m.iconClose.Free()
		m.iconCloseHighlight.Free()
		m.icon.Free()
	})
	return m
}

func (m *TButton) iconChange(sender lcl.IObject) {
	if m.IsDisable {
		return
	}
	m.Invalidate()
}

func (m *TButton) enter(sender lcl.IObject) {
	if m.IsDisable {
		return
	}
	m.isEnter = true
	m.Invalidate()
	if m.onMouseEnter != nil {
		m.onMouseEnter(sender)
	}
}

func (m *TButton) isCloseArea(X int32, Y int32) bool {
	rect := m.ClientRect()
	closeX := rect.Width() - m.iconClose.Width() - 10
	closeY := rect.Height()/2 - m.iconClose.Height()/2
	return X >= closeX && X <= rect.Width()-10 && Y >= closeY && Y <= rect.Height()/2+m.iconClose.Height()
}

func (m *TButton) move(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
	lcl.Screen.SetCursor(types.CrDefault)
	if m.IsDisable {
		return
	}
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

func (m *TButton) leave(sender lcl.IObject) {
	if m.IsDisable {
		return
	}
	m.isEnter = false
	m.isEnterClose = false
	m.Invalidate()
	if m.onMouseLeave != nil {
		m.onMouseLeave(sender)
	}
}

func (m *TButton) down(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.IsDisable {
		return
	}
	if m.isCloseArea(X, Y) {
		if m.onCloseClick != nil {
			m.onCloseClick(sender)
		}
	} else {
		m.isDown = true
		m.Invalidate()
		if m.onMouseDown != nil {
			m.onMouseDown(sender, button, shift, X, Y)
		}
	}
}

func (m *TButton) up(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.IsDisable {
		return
	}
	m.isDown = false
	m.Invalidate()
	if m.onMouseUp != nil {
		m.onMouseUp(sender, button, shift, X, Y)
	}
}

func (m *TButton) drawRoundedGradientButton(canvas lcl.ICanvas, rect types.TRect) {
	text := m.Caption()
	// 创建图像对象
	img := lcl.NewLazIntfImageWithIntX2RawImageQueryFlags(rect.Width(), rect.Height(), types.NewSet(types.RiqfRGB, types.RiqfAlpha))
	defer img.Free()

	startColor := m.startColor
	endColor := m.endColor

	if m.isEnter {
		startColor = darkenColor(startColor, 0.1)
		endColor = darkenColor(endColor, 0.1)
	}
	if m.isDown {
		startColor = darkenColor(startColor, 0.2)
		endColor = darkenColor(endColor, 0.2)
	}

	// 获取起始颜色分量
	startR := colors.Red(startColor)
	startG := colors.Green(startColor)
	startB := colors.Blue(startColor)

	// 获取结束颜色分量
	endR := colors.Red(endColor)
	endG := colors.Green(endColor)
	endB := colors.Blue(endColor)

	// 创建垂直渐变（带抗锯齿圆角）
	imgHeight := img.Height()
	imgWidth := img.Width()
	for y := 0; y < int(imgHeight); y++ {
		ratio := float64(y) / float64(imgHeight-1)
		r := round(float64(startR)*(1-ratio) + float64(endR)*ratio)
		g := round(float64(startG)*(1-ratio) + float64(endG)*ratio)
		b := round(float64(startB)*(1-ratio) + float64(endB)*ratio)
		curColor := lcl.TFPColor{Red: uint16(r) << 8, Green: uint16(g) << 8, Blue: uint16(b) << 8}
		// 注意：Alpha会在内循环中为每个像素单独设置
		for x := 0; x < int(imgWidth); x++ {
			alphaFactor := m.calculateRoundedAlpha(int32(x), int32(y), imgWidth, imgHeight, m.radius)
			actualAlpha := round(float64(m.alpha) * float64(alphaFactor))
			curColor.Alpha = uint16(actualAlpha) << 8
			img.SetColors(int32(x), int32(y), curColor)
		}
	}
	// 创建临时位图并加载图像数据
	tempBMap := lcl.NewBitmap()
	tempBMap.SetSize(rect.Width(), rect.Height())
	tempBMap.SetPixelFormat(types.Pf32bit)
	tempBMap.LoadFromIntfImage(img)
	// 绘制到目标画布
	canvas.DrawWithIntX2Graphic(rect.Left, rect.Top, tempBMap)
	defer tempBMap.Free()

	// 绘制按钮文字（在原始画布上绘制，确保文字不透明）
	brush := canvas.BrushToBrush()
	brush.SetStyle(types.BsClear)

	// 计算文字位置
	textSize := canvas.TextExtentWithUnicodestring(text)
	textX := rect.Left + m.TextOffSetX + (rect.Width()-textSize.Cx)/2
	textY := rect.Top + m.TextOffSetY + (rect.Height()-textSize.Cy)/2

	// 计算文字宽度截取
	if len(text) > 0 && m.iconFavorite.Width() > 0 { // 当有文本时才计算截取
		//textWidth := canvas.GetTextWidthWithUnicodestring(text)
		//fmt.Println("text:", text, textWidth, textWidth/int32(len(text)), m.iconFavorite.Width())
		leftIconSize := int32(10 + 20) // 边侧图标的距离， left 10, icon 宽高 20
		if textX <= leftIconSize {
			textX = leftIconSize
		}

	}

	// 绘制文字阴影（增强可读性）
	//canvas.FontToFont().SetColor(colors.ClBlack)
	//canvas.TextOutWithIntX2Unicodestring(textX+1, textY+1, text)

	// 绘制主文字
	//canvas.FontToFont().SetColor(colors.ClWhite)
	canvas.TextOutWithIntX2Unicodestring(textX, textY, text)

	// 绘制图标 favorite 在左
	favY := rect.Height()/2 - m.iconFavorite.Height()/2
	canvas.DrawWithIntX2Graphic(10, favY, m.iconFavorite.Graphic())
	// 绘制图标 close 在右
	iconClose := m.iconClose
	if m.isEnterClose {
		iconClose = m.iconCloseHighlight
	}
	closeX := rect.Width() - iconClose.Width() - 10
	closeY := rect.Height()/2 - iconClose.Height()/2
	canvas.DrawWithIntX2Graphic(closeX, closeY, iconClose.Graphic())

	// 绘制图标 icon, 在中间位置
	iconW, iconH := m.icon.Width(), m.icon.Height()
	iconX := rect.Left + (rect.Width()-iconW)/2
	iconY := rect.Top + (rect.Height()-iconH)/2
	canvas.DrawWithIntX2Graphic(iconX, iconY, m.icon.Graphic())
}

func (m *TButton) SetIcon(filePath string) {
	m.icon.LoadFromFile(filePath)
}

func (m *TButton) SetIconFavorite(filePath string) {
	m.iconFavorite.LoadFromFile(filePath)
}

func (m *TButton) SetIconClose(filePath string) {
	path, name := filepath.Split(filePath)
	m.iconClose.LoadFromFile(filePath)
	ns := strings.Split(name, ".")
	m.iconCloseHighlight.LoadFromFile(filepath.Join(path, ns[0]+"_enter.png"))
}

func (m *TButton) paint(sender lcl.IObject) {
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

func (m *TButton) SetStartColor(color colors.TColor) {
	m.startColor = color
}

func (m *TButton) SetEndColor(color colors.TColor) {
	m.endColor = color
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

// 计算圆角矩形中某点的抗锯齿透明度因子 (0.0 ~ 1.0)
func (m *TButton) calculateRoundedAlpha(x, y, width, height, radius int32) float32 {
	// 计算实际可用最大半径（不超过尺寸限制）
	maxRadius := min(width/2, height/2)
	if radius > maxRadius {
		radius = maxRadius
	}
	// 如果半径被限制为0，直接返回不透明
	if radius <= 0 {
		return 1.0
	}

	var (
		cornerX, cornerY int32
		d                float32
		inCorner         bool
	)

	// 左上角区域
	if m.RoundedCorner.In(RcLeftTop) && x < radius && y < radius {
		cornerX = radius
		cornerY = radius
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
		inCorner = true
	} else if m.RoundedCorner.In(RcRightTop) && x >= width-radius && y < radius {
		// 右上角区域
		cornerX = width - radius - 1
		cornerY = radius
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
		inCorner = true
	} else if m.RoundedCorner.In(RcLeftBottom) && x < radius && y >= height-radius {
		// 左下角区域
		cornerX = radius
		cornerY = height - radius - 1
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
		inCorner = true
	} else if m.RoundedCorner.In(RcRightBottom) && x >= width-radius && y >= height-radius {
		// 右下角区域
		cornerX = width - radius - 1
		cornerY = height - radius - 1
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
		inCorner = true
	}

	if !inCorner {
		// 非圆角区域：检查是否在有效矩形内
		if x >= radius && x < width-radius &&
			y >= radius && y < height-radius {
			return 1.0 // 中央矩形区域
		}
		// 边缘非圆角区域
		return 1.0
	}

	// 抗锯齿过渡区域（像素宽度）
	const transition = 1.0
	innerRadius := float32(radius) - transition

	// 完全在圆角内
	if d <= innerRadius {
		return 1.0
	}
	// 完全在圆角外
	if d >= float32(radius)+transition {
		return 0.0
	}

	// 在过渡区域内（平滑渐变）
	return 1.0 - (d-innerRadius)/(2*transition)
}

// 辅助函数：整数最小值
func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func darkenColor(color types.TColor, factor float64) types.TColor {
	R := colors.Red(color)
	G := colors.Green(color)
	B := colors.Blue(color)

	R = byte(round(float64(R) * (1.0 - factor)))
	G = byte(round(float64(G) * (1.0 - factor)))
	B = byte(round(float64(B) * (1.0 - factor)))
	return colors.RGBToColor(R, G, B)
}

func round(v float64) float64 {
	return math.Round(v)
}

func sqr(x int32) int32 {
	return x * x
}

func sqrt(v float64) float32 {
	return float32(math.Sqrt(v))
}
