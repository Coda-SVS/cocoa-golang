package main

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/config"
	_ "github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/window"
)

var (
	isExit bool = false
)

func main() {
	defer audio.Dispose()
	defer config.WriteConfig()

	wnd := g.NewMasterWindow("COCOA", 1400, 800, 0)
	imgui.StyleColorsDark()

	wnd.SetCloseCallback(callbackClose)
	wnd.Run(window.MainWindowGUILoop)
}

func callbackClose() bool {
	isExit = true
	return isExit
}
