package main

import (
	"runtime"
	"time"

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

	frameRefresherQuitChan chan any
)

func main() {
	runtime.LockOSThread()

	imguiw.InitImgui("COCOA", 1400, 800)
	imgui.StyleColorsDark()
	imguiw.SetBeforeDestroyContextCallback(func() {
		if frameRefresherQuitChan != nil {
			close(frameRefresherQuitChan)
		}

		audio.Dispose()
		dsp.Dispose()
		config.WriteConfig()
		log.Dispose()
	})

	mainWindow = window.NewMainWindow()

	go frameRefresher()

	imguiw.Run(log.PanicLogHandler(log.RootLogger(), util.PanicToErrorW(mainWindowGUILoop)))
}

func mainWindowGUILoop() {
	mainWindow.View()
}

func frameRefresher() {
	msgChan := audio.AudioStreamBroker().Subscribe()

	for msg := range msgChan {
		switch msg {
		case audio.EnumAudioStreamStarted:
			if frameRefresherQuitChan == nil {
				imguiw.Context.WaitGroup().Add(1)
				frameRefresherQuitChan = make(chan any)
				go func() {
					backend := imguiw.Context.Backend()
					ticker := time.NewTicker(time.Second / 60)
					for {
						select {
						case <-ticker.C:
							backend.Refresh()
						case <-frameRefresherQuitChan:
							imguiw.Context.WaitGroup().Done()
							return
						}
					}
				}()
			}

		case audio.EnumAudioStreamPaused:
			fallthrough
		case audio.EnumAudioStreamStoped:
			if frameRefresherQuitChan != nil {
				close(frameRefresherQuitChan)
				frameRefresherQuitChan = nil
			}
		}
	}
}
