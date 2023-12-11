package imguiw

import imgui "github.com/AllenDang/cimgui-go"

var (
	windowClass *imgui.WindowClass
)

func init() {
	windowClass = imgui.NewWindowClass()
	windowClass.SetViewportFlagsOverrideSet(imgui.ViewportFlagsNoAutoMerge)
}

func ApplySubWindowClass() {
	imgui.SetNextWindowClass(windowClass)
}
