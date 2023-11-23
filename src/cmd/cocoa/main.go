package main

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	_ "github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/window"
)

func main() {
	wnd := g.NewMasterWindow("COCOA", 1400, 800, 0)
	imgui.StyleColorsDark()
	wnd.Run(window.MainWindowGUILoop)
}
