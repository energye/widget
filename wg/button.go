package wg

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"math"
)

type RoundedCorner = int32

const (
	RcLeftTop RoundedCorner = iota
	RcRightTop
	RcLeftBottom
	RcRightBottom
)

type RoundedCorners = types.TSet

type TButtonIcon struct {
	Icon []byte
	X    int32
	Y    int32
}

type TButton struct {
	lcl.ICustomGraphicControl
	startColor    colors.TColor  // 按钮起始渐变颜色
	endColor      colors.TColor  // 按钮结束渐变颜色
	activeColor   float32        // 按钮激活颜色深度 0.0 ~ 1.0
	alpha         byte           // 透明度 0 ~ 255
	radius        int32          // 圆角度
	isEnter       bool           // 鼠标是否移入
	isDown        bool           // 鼠标是否按下
	RoundedCorner RoundedCorners // 按钮圆角方向，默认四角
	Icons         []TButtonIcon  // 按钮上的图标
	// 事件
	onPaint      lcl.TNotifyEvent
	onMouseEnter lcl.TNotifyEvent
	onMouseLeave lcl.TNotifyEvent
	onMouseDown  lcl.TMouseEvent
	onMouseUp    lcl.TMouseEvent
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
	m.radius = 15
	m.ICustomGraphicControl.SetOnPaint(m.paint)
	m.ICustomGraphicControl.SetOnMouseEnter(m.enter)
	m.ICustomGraphicControl.SetOnMouseLeave(m.leave)
	m.ICustomGraphicControl.SetOnMouseDown(m.down)
	m.ICustomGraphicControl.SetOnMouseUp(m.up)
	m.RoundedCorner = types.NewSet(RcLeftTop, RcRightTop, RcLeftBottom, RcRightBottom)
	return m
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

func (m *TButton) enter(sender lcl.IObject) {
	m.isEnter = true
	m.Invalidate()
	if m.onMouseEnter != nil {
		m.onMouseEnter(sender)
	}
}

func (m *TButton) leave(sender lcl.IObject) {
	m.isEnter = false
	m.Invalidate()
	if m.onMouseLeave != nil {
		m.onMouseLeave(sender)
	}
}

func (m *TButton) down(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	m.isDown = true
	m.Invalidate()
	if m.onMouseDown != nil {
		m.onMouseDown(sender, button, shift, X, Y)
	}
}

func (m *TButton) up(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
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
	textX := rect.Left + (rect.Width()-textSize.Cx)/2
	textY := rect.Top + (rect.Height()-textSize.Cy)/2

	// 计算文字宽度截取
	//textWidth := canvas.GetTextWidthWithUnicodestring(text)
	//fmt.Println("text:", text, textWidth, textWidth/int32(len(text)))

	// 绘制文字阴影（增强可读性）
	//canvas.FontToFont().SetColor(colors.ClBlack)
	//canvas.TextOutWithIntX2Unicodestring(textX+1, textY+1, text)

	// 绘制主文字
	//canvas.FontToFont().SetColor(colors.ClWhite)
	canvas.TextOutWithIntX2Unicodestring(textX, textY, text)
}

func (m *TButton) paint(sender lcl.IObject) {
	m.drawRoundedGradientButton(m.Canvas(), m.ClientRect())
	if m.onPaint != nil {
		m.onPaint(sender)
	}
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

// 计算圆角矩形中某点的抗锯齿透明度因子 (0.0 ~ 1.0)
func (m *TButton) calculateRoundedAlpha(x, y, width, height, radius int32) float32 {
	var (
		cornerX, cornerY int32
		d                float32
	)
	// 左上角区域
	if m.RoundedCorner.In(RcLeftTop) && x < radius && y < radius {
		cornerX = radius
		cornerY = radius
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
	} else if m.RoundedCorner.In(RcRightTop) && x >= width-radius && y < radius {
		// 右上角区域
		cornerX = width - radius - 1
		cornerY = radius
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
	} else if m.RoundedCorner.In(RcLeftBottom) && x < radius && y >= height-radius {
		// 左下角区域
		cornerX = radius
		cornerY = height - radius - 1
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
	} else if m.RoundedCorner.In(RcRightBottom) && x >= width-radius && y >= height-radius {
		// 右下角区域
		cornerX = width - radius - 1
		cornerY = height - radius - 1
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
	} else {
		return 1.0
	}
	// 计算圆角区域的抗锯齿透明度
	if d < float32(radius)-1.0 {
		return 1.0 // 完全在圆角内
	} else if d > float32(radius)+1.0 {
		return 0.0 // 完全在圆角外
	} else {
		// 在过渡区域（1像素宽度），线性插值
		return 1.0 - (d-(float32(radius)-1.0))/2.0
	}
}
