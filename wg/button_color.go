package wg

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// TButtonColor 按钮颜色
type TButtonColor struct {
	start  colors.TColor     // 按钮起始渐变颜色
	end    colors.TColor     // 按钮结束渐变颜色
	img    lcl.ILazIntfImage //
	bitMap lcl.IBitmap       //
}

func NewButtonColor(start, end colors.TColor) *TButtonColor {
	m := &TButtonColor{
		start:  start,
		end:    end,
		img:    lcl.NewLazIntfImageWithIntX2RawImageQueryFlags(0, 0, types.NewSet(types.RiqfRGB, types.RiqfAlpha)),
		bitMap: lcl.NewBitmap(),
	}
	m.bitMap.SetPixelFormat(types.Pf32bit)
	return m
}

// 强制执行绘制
func (m *TButtonColor) forcePaint(roundedCorners TRoundedCorners, rect types.TRect, alpha byte, radius int32) {
	if m.img.Width() != rect.Width() || m.img.Height() != rect.Height() {
		m.img.SetSize(rect.Width(), rect.Height())
	}
	if m.bitMap.Width() != rect.Width() || m.bitMap.Height() != rect.Height() {
		m.bitMap.SetSize(rect.Width(), rect.Height())
	}
	m.doPaint(roundedCorners, alpha, radius)
}

// 绘制根据TRect变化
func (m *TButtonColor) paint(roundedCorners TRoundedCorners, rect types.TRect, alpha byte, radius int32) {
	isPaint := false
	if m.img.Width() != rect.Width() || m.img.Height() != rect.Height() {
		m.img.SetSize(rect.Width(), rect.Height())
		isPaint = true
	}
	if m.bitMap.Width() != rect.Width() || m.bitMap.Height() != rect.Height() {
		m.bitMap.SetSize(rect.Width(), rect.Height())
		isPaint = true
	}

	if !isPaint {
		return
	}
	m.doPaint(roundedCorners, alpha, radius)
}

// 执行绘制
func (m *TButtonColor) doPaint(roundedCorners TRoundedCorners, alpha byte, radius int32) {
	// 获取起始颜色分量
	startR := colors.Red(m.start)
	startG := colors.Green(m.start)
	startB := colors.Blue(m.start)
	// 获取结束颜色分量
	endR := colors.Red(m.end)
	endG := colors.Green(m.end)
	endB := colors.Blue(m.end)
	// 处理垂直渐变（带抗锯齿圆角）
	imgHeight := m.img.Height()
	imgWidth := m.img.Width()
	for y := 0; y < int(imgHeight); y++ {
		ratio := float64(y) / float64(imgHeight-1)
		r := round(float64(startR)*(1-ratio) + float64(endR)*ratio)
		g := round(float64(startG)*(1-ratio) + float64(endG)*ratio)
		b := round(float64(startB)*(1-ratio) + float64(endB)*ratio)
		curColor := lcl.TFPColor{Red: uint16(r) << 8, Green: uint16(g) << 8, Blue: uint16(b) << 8}
		// 注意：Alpha会在内循环中为每个像素单独设置
		for x := 0; x < int(imgWidth); x++ {
			alphaFactor := m.calculateRoundedAlpha(roundedCorners, int32(x), int32(y), imgWidth, imgHeight, radius)
			actualAlpha := round(float64(alpha) * float64(alphaFactor))
			curColor.Alpha = uint16(actualAlpha) << 8
			m.img.SetColors(int32(x), int32(y), curColor)
		}
	}
	// 位图加载图像数据
	m.bitMap.LoadFromIntfImage(m.img)
}

// 计算圆角矩形中某点的抗锯齿透明度因子 (0.0 ~ 1.0)
func (m *TButtonColor) calculateRoundedAlpha(roundedCorners TRoundedCorners, x, y, width, height, radius int32) float32 {
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
	if roundedCorners.In(RcLeftTop) && x < radius && y < radius {
		cornerX = radius
		cornerY = radius
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
		inCorner = true
	} else if roundedCorners.In(RcRightTop) && x >= width-radius && y < radius {
		// 右上角区域
		cornerX = width - radius - 1
		cornerY = radius
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
		inCorner = true
	} else if roundedCorners.In(RcLeftBottom) && x < radius && y >= height-radius {
		// 左下角区域
		cornerX = radius
		cornerY = height - radius - 1
		d = sqrt(float64(sqr(x-cornerX) + sqr(y-cornerY)))
		inCorner = true
	} else if roundedCorners.In(RcRightBottom) && x >= width-radius && y >= height-radius {
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
