package wg

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"strconv"
)

var (
	activeColor   = colors.RGBToColor(60, 70, 80)
	defaultColor  = colors.RGBToColor(86, 88, 100)
	defaultPrefix = "Tab"
	defaultHeight = int32(25)
)

type TTab struct {
	lcl.ICustomPanel          //
	pages            []*TPage // 页列表
	widths           int32    // 页签总宽度
	deleting         bool     // 正在删除中的 page
	scrollLeftBtn    *TButton // tab 滚动导航按钮 左滚动
	scrollRightBtn   *TButton // tab 滚动导航按钮 右滚动
}

type TPage struct {
	lcl.ICustomPanel
	active bool     // 是否激活
	show   bool     // 是否显示
	tab    *TTab    // 所属的tab
	button *TButton // 按钮
}

func NewTab(owner lcl.IComponent) *TTab {
	tab := &TTab{}
	tab.ICustomPanel = lcl.NewCustomPanel(owner)
	tab.SetBevelInner(types.BvNone)
	tab.SetBevelOuter(types.BvNone)
	//tab.SetColor(colors.ClRed)
	tab.SetBorderStyleToBorderStyle(types.BsNone)
	tab.initScrollBtn()
	return tab
}

// 初始化滚动导航按钮
func (m *TTab) initScrollBtn() {
	m.scrollLeftBtn = NewButton(m)
	m.scrollRightBtn = NewButton(m)

}

func (m *TTab) NewPage() *TPage {
	page := new(TPage)
	page.tab = m
	button := NewButton(m)
	button.SetAutoSize(true)
	button.SetShowHint(true)
	button.SetText(defaultPrefix + strconv.Itoa(len(m.pages)))
	button.Font().SetSize(9)
	button.Font().SetColor(colors.Cl3DFace)
	button.SetStartColor(defaultColor)
	button.SetEndColor(defaultColor)
	button.RoundedCorner = button.RoundedCorner.Exclude(RcLeftBottom).Exclude(RcRightBottom)
	button.SetRadius(0)
	button.SetAlpha(255)
	button.SetHeight(defaultHeight)
	button.SetParent(m)
	page.button = button

	sheet := lcl.NewPanel(m)
	sheet.SetBevelInner(types.BvNone)
	sheet.SetBevelOuter(types.BvNone)
	sheet.SetBorderStyleToBorderStyle(types.BsNone)
	tabRect := m.ClientRect()
	sheet.SetTop(button.Height())
	sheet.SetHeight(tabRect.Height() - button.Height())
	sheet.SetWidth(tabRect.Width())
	sheet.SetAlign(types.AlCustom)
	sheet.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
	sheet.SetParent(m)
	page.ICustomPanel = sheet

	m.pages = append(m.pages, page)
	m.recalculatePosition()
	page.initEvent()
	m.HiddenAllActivated()
	page.Active(true)
	return page
}

// 重新计算位置, 在隐藏/移除时使用
func (m *TTab) recalculatePosition() {
	var widths int32
	for _, page := range m.pages {
		br := page.button.BoundsRect()
		width := br.Width()
		br.Left = widths
		br.SetWidth(width)
		page.button.SetBoundsRect(br)
		widths += br.Width()
	}
	m.widths = widths
}

// HiddenAllActivated 隐藏所有激活页面
func (m *TTab) HiddenAllActivated() {
	for _, page := range m.pages {
		if page.active {
			page.Active(false)
		}
	}
}

//func (m *TTab)

// 删除指定 page
func (m *TTab) RemovePage(removePage *TPage) {
	removeIndex := -1 // 存放当前删除page的索引
	for i, page := range m.pages {
		if page == removePage {
			removeIndex = i
			m.pages = append(m.pages[:i], m.pages[i+1:]...)
			break
		}
	}
	// 重新计算 button 位置
	m.recalculatePosition()
	if removePage.active {
		// 根据删除索引获取要显示的 page
		var showPage *TPage
		if removeIndex != -1 && removeIndex < len(m.pages) {
			showPage = m.pages[removeIndex] // 显示当前索引的 page, 也就是删除后的下一个
		} else if len(m.pages) > 0 {
			showPage = m.pages[0] // 显示第一个 page
		}
		if showPage != nil {
			m.HiddenAllActivated()
			showPage.Active(true)
		}
	}
	m.deleting = false
}

// 删除掉自己
func (m *TPage) Remove() {
	m.button.SetOnClick(nil)
	m.button.SetOnCloseClick(nil)
	// 先隐藏掉
	m.button.Hide()
	m.ICustomPanel.Hide()
	// 在page里删除自己
	m.tab.RemovePage(m)
	// 最后释放掉
	m.button.Free()
	m.ICustomPanel.Free()
	m.tab = nil
}

// 激活自己, 会取消激活其它的
func (m *TPage) Active(active bool) {
	m.active = active
	if active {
		m.ICustomPanel.Show()
		m.button.SetStartColor(activeColor)
		m.button.SetEndColor(activeColor)
		m.button.Invalidate()
	} else {
		m.ICustomPanel.Hide()
		m.button.SetStartColor(defaultColor)
		m.button.SetEndColor(defaultColor)
		m.button.Invalidate()
	}
}

// 隐藏自己, button 和 page 同时隐藏
func (m *TPage) Hide() {

}

// 显示自己, button 和 page 同时显示
func (m *TPage) Show() {

}

func (m *TPage) initEvent() {
	m.button.SetOnClick(func(sender lcl.IObject) {
		m.tab.HiddenAllActivated()
		m.Active(true)
	})
	m.button.SetOnCloseClick(func(sender lcl.IObject) {
		if m.tab.deleting {
			return
		}
		m.tab.deleting = true
		lcl.RunOnMainThreadAsync(func(id uint32) {
			m.Remove()
		})
	})
	m.SetOnResize(func(sender lcl.IObject) {
		br := m.BoundsRect()
		println("Width:", br.Width())
	})
}

func (m *TPage) Button() *TButton {
	return m.button
}
