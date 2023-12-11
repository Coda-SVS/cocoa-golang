package imguiw

import (
	"math"

	imgui "github.com/AllenDang/cimgui-go"
)

type imguiStyleState struct {
	WindowPadding             imgui.Vec2
	WindowRounding            float32
	WindowMinSize             imgui.Vec2
	ChildRounding             float32
	PopupRounding             float32
	FramePadding              imgui.Vec2
	FrameRounding             float32
	ItemSpacing               imgui.Vec2
	ItemInnerSpacing          imgui.Vec2
	CellPadding               imgui.Vec2
	TouchExtraPadding         imgui.Vec2
	IndentSpacing             float32
	ColumnsMinSpacing         float32
	ScrollbarSize             float32
	ScrollbarRounding         float32
	GrabMinSize               float32
	GrabRounding              float32
	LogSliderDeadzone         float32
	TabRounding               float32
	TabMinWidthForCloseButton float32
	SeparatorTextPadding      imgui.Vec2
	DockingSeparatorSize      float32
	DisplayWindowPadding      imgui.Vec2
	DisplaySafeAreaPadding    imgui.Vec2
	MouseCursorScale          float32
}

type ImguiDPI struct {
	currentStyle   *imgui.Style
	currentContext *imgui.Context
	currentIO      *imgui.IO
	currentDPI     float32
	currentState   *imguiStyleState
}

func NewImguiDPI(style *imgui.Style, context *imgui.Context, io *imgui.IO) *ImguiDPI {
	if style == nil {
		style = imgui.CurrentStyle()
	}
	if context == nil {
		context = imgui.CurrentContext()
	}
	if io == nil {
		io = imgui.CurrentIO()
	}

	imguiDPI := &ImguiDPI{
		currentStyle:   style,
		currentContext: context,
		currentIO:      io,
	}
	return imguiDPI
}

func (d *ImguiDPI) GetDPI() float32 {
	return d.currentDPI
}

func (d *ImguiDPI) UpdateDPI() {
	if d.currentState == nil {
		d.backupCurrentStyle()
	}

	dpiScale := d.currentContext.CurrentDpiScale()

	if d.currentDPI != dpiScale {
		d.currentIO.SetFontGlobalScale(dpiScale)
		d.applyStyle(dpiScale)
		d.currentDPI = dpiScale
	}
}

func (d *ImguiDPI) applyStyle(dpiScale float32) {
	styleHandler := d.currentStyle
	currentState := d.currentState
	styleHandler.SetWindowPadding(currentState.WindowPadding.Mul(dpiScale))
	styleHandler.SetWindowRounding(currentState.WindowRounding * dpiScale)
	styleHandler.SetWindowMinSize(currentState.WindowMinSize.Mul(dpiScale))
	styleHandler.SetChildRounding(currentState.ChildRounding * dpiScale)
	styleHandler.SetPopupRounding(currentState.PopupRounding * dpiScale)
	styleHandler.SetFramePadding(currentState.FramePadding.Mul(dpiScale))
	styleHandler.SetFrameRounding(currentState.FrameRounding * dpiScale)
	styleHandler.SetItemSpacing(currentState.ItemSpacing.Mul(dpiScale))
	styleHandler.SetItemInnerSpacing(currentState.ItemInnerSpacing.Mul(dpiScale))
	styleHandler.SetCellPadding(currentState.CellPadding.Mul(dpiScale))
	styleHandler.SetTouchExtraPadding(currentState.TouchExtraPadding.Mul(dpiScale))
	styleHandler.SetIndentSpacing(currentState.IndentSpacing * dpiScale)
	styleHandler.SetColumnsMinSpacing(currentState.ColumnsMinSpacing * dpiScale)
	styleHandler.SetScrollbarSize(currentState.ScrollbarSize * dpiScale)
	styleHandler.SetScrollbarRounding(currentState.ScrollbarRounding * dpiScale)
	styleHandler.SetGrabMinSize(currentState.GrabMinSize * dpiScale)
	styleHandler.SetGrabRounding(currentState.GrabRounding * dpiScale)
	styleHandler.SetLogSliderDeadzone(currentState.LogSliderDeadzone * dpiScale)
	styleHandler.SetTabRounding(currentState.TabRounding * dpiScale)
	if math.MaxFloat32 != currentState.TabMinWidthForCloseButton {
		styleHandler.SetTabMinWidthForCloseButton(currentState.TabMinWidthForCloseButton * dpiScale)
	} else {
		styleHandler.SetTabMinWidthForCloseButton(math.MaxFloat32)
	}
	styleHandler.SetSeparatorTextPadding(currentState.SeparatorTextPadding.Mul(dpiScale))
	styleHandler.SetDockingSeparatorSize(currentState.DockingSeparatorSize * dpiScale)
	styleHandler.SetDisplayWindowPadding(currentState.DisplayWindowPadding.Mul(dpiScale))
	styleHandler.SetDisplaySafeAreaPadding(currentState.DisplaySafeAreaPadding.Mul(dpiScale))
	styleHandler.SetMouseCursorScale(currentState.MouseCursorScale * dpiScale)
}

func (d *ImguiDPI) backupCurrentStyle() {
	styleHandler := d.currentStyle
	newState := imguiStyleState{}

	newState.WindowPadding = styleHandler.WindowPadding()
	newState.WindowRounding = styleHandler.WindowRounding()
	newState.WindowMinSize = styleHandler.WindowMinSize()
	newState.ChildRounding = styleHandler.ChildRounding()
	newState.PopupRounding = styleHandler.PopupRounding()
	newState.FramePadding = styleHandler.FramePadding()
	newState.FrameRounding = styleHandler.FrameRounding()
	newState.ItemSpacing = styleHandler.ItemSpacing()
	newState.ItemInnerSpacing = styleHandler.ItemInnerSpacing()
	newState.CellPadding = styleHandler.CellPadding()
	newState.TouchExtraPadding = styleHandler.TouchExtraPadding()
	newState.IndentSpacing = styleHandler.IndentSpacing()
	newState.ColumnsMinSpacing = styleHandler.ColumnsMinSpacing()
	newState.ScrollbarSize = styleHandler.ScrollbarSize()
	newState.ScrollbarRounding = styleHandler.ScrollbarRounding()
	newState.GrabMinSize = styleHandler.GrabMinSize()
	newState.GrabRounding = styleHandler.GrabRounding()
	newState.LogSliderDeadzone = styleHandler.LogSliderDeadzone()
	newState.TabRounding = styleHandler.TabRounding()
	newState.TabMinWidthForCloseButton = styleHandler.TabMinWidthForCloseButton()
	newState.SeparatorTextPadding = styleHandler.SeparatorTextPadding()
	newState.DockingSeparatorSize = styleHandler.DockingSeparatorSize()
	newState.DisplayWindowPadding = styleHandler.DisplayWindowPadding()
	newState.DisplaySafeAreaPadding = styleHandler.DisplaySafeAreaPadding()
	newState.MouseCursorScale = styleHandler.MouseCursorScale()

	d.currentState = &newState
}
