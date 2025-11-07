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

// forcePaint 强制绘制按钮颜色
// roundedCorners: 圆角设置
// rect: 绘制区域矩形
// alpha: 透明度值
// radius: 圆角半径
func (m *TButtonColor) forcePaint(roundedCorners TRoundedCorners, rect types.TRect, alpha byte, radius int32) {
	if m.img.Width() != rect.Width() || m.img.Height() != rect.Height() {
		m.img.SetSize(rect.Width(), rect.Height())
	}
	if m.bitMap.Width() != rect.Width() || m.bitMap.Height() != rect.Height() {
		m.bitMap.SetSize(rect.Width(), rect.Height())
	}
	m.doPaint(roundedCorners, alpha, radius)
}

// paint 绘制按钮颜色
// roundedCorners: 圆角设置
// rect: 绘制区域矩形
// alpha: 透明度值
// radius: 圆角半径
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

// doPaint 绘制带有圆角和透明度的垂直渐变按钮图像。
// 参数:
//
//	roundedCorners: 指定哪些角落需要绘制为圆角
//	alpha: 图像的整体透明度，取值范围 0-255
//	radius: 圆角的半径大小
func (m *TButtonColor) doPaint(roundedCorners TRoundedCorners, alpha byte, radius int32) {
	// 提取起始颜色和结束颜色的 RGB 分量，用于计算渐变过程中的颜色插值
	startR := colors.Red(m.start)
	startG := colors.Green(m.start)
	startB := colors.Blue(m.start)
	// 获取结束颜色分量
	endR := colors.Red(m.end)
	endG := colors.Green(m.end)
	endB := colors.Blue(m.end)
	// 处理垂直渐变（带抗锯齿圆角）
	// 遍历图像每一行，根据当前行位置计算颜色渐变比例，并逐像素设置颜色与透明度
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
	// 将处理好的图像数据加载到位图对象中，供后续使用
	m.bitMap.LoadFromIntfImage(m.img)
}

// calculateRoundedAlpha 根据给定的圆角信息计算按钮颜色在指定坐标点 (x, y) 处的 alpha 值，
// 用于实现抗锯齿效果的圆角矩形绘制。
//
// 参数说明：
//
//	roundedCorners: 指定哪些角落需要绘制为圆角（TRoundedCorners 类型）
//	x: 当前像素点的横坐标
//	y: 当前像素点的纵坐标
//	width: 矩形区域的宽度
//	height: 矩形区域的高度
//	radius: 圆角半径
//
// 返回值：
//
//	float32: 当前点的 alpha 值，范围 [0.0, 1.0]，表示透明度（0 为完全透明，1 为不透明）
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
		cornerX, cornerY int32   // 圆角顶点坐标
		d                float32 // 距离圆心的距离
		inCorner         bool    // 是否在圆角内
	)
	// 判断当前点是否位于某个圆角区域内，并计算其到对应圆心的距离
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
		if x >= radius && x < width-radius && y >= radius && y < height-radius {
			return 1.0 // 中央矩形区域
		}
		// 边缘非圆角区域
		return 1.0
	}
	// 抗锯齿过渡处理：根据距离决定 alpha 渐变值
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
