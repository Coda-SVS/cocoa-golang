package main

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/config"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
	"github.com/Kor-SVS/cocoa/src/ui/window"
	"github.com/Kor-SVS/cocoa/src/util"
)

var (
	mainWindow imguiw.Widget
)

func main() {
	defer audio.Dispose()
	defer config.WriteConfig()

	imguiw.InitImgui("COCOA", 1400, 800)
	imgui.StyleColorsDark()

	mainWindow = window.NewMainWindow()

	imguiw.Run(log.PanicLogHandler(log.RootLogger(), util.PanicToErrorW(mainWindowGUILoop)))
}

func mainWindowGUILoop() {
	mainWindow.View()
}
