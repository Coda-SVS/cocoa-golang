package window

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
	"github.com/Kor-SVS/cocoa/src/ui/widget"
	"github.com/sqweek/dialog"
)

type MainWindow struct {
	title string
	State MainWindowState
}

type MainWindowState struct {
	IsOpen           bool
	LeftSidePanelPos float32
}

func NewMainWindow() (window *MainWindow) {
	window = &MainWindow{}
	window.title = "Main Window"
	window.State = MainWindowState{
		LeftSidePanelPos: 400,
	}
	return window
}

func (mw *MainWindow) Title() string {
	mtx := imguiw.Context.Mutex()
	mtx.Lock()
	defer mtx.Unlock()

	return mw.title
}

func (mw *MainWindow) IsOpen() bool {
	mtx := imguiw.Context.Mutex()
	mtx.Lock()
	defer mtx.Unlock()

	return mw.State.IsOpen
}

func (mw *MainWindow) SetIsOpen(value bool) {
	mtx := imguiw.Context.Mutex()
	mtx.Lock()
	defer mtx.Unlock()

	mw.State.IsOpen = value
}

func (mw *MainWindow) View() {
	backend := imguiw.Context.Backend()

	pos := imgui.MainViewport().Pos()
	sizeX, sizeY := backend.DisplaySize()

	imgui.SetNextWindowPos(pos)
	imgui.SetNextWindowSize(imgui.Vec2{X: float32(sizeX), Y: float32(sizeY)})
	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 0)
	if imgui.BeginV(imguiw.T(mw.title), &mw.State.IsOpen,
		imgui.WindowFlagsNoDocking|
			imgui.WindowFlagsNoTitleBar|
			imgui.WindowFlagsNoCollapse|
			imgui.WindowFlagsNoScrollbar|
			imgui.WindowFlagsNoMove|
			imgui.WindowFlagsNoResize|
			imgui.WindowFlagsMenuBar) {
		if imgui.BeginMenuBar() {
			if imgui.BeginMenu(imguiw.T("File")) {
				if imgui.MenuItemBool(imguiw.T("OpenFile")) {
					openFile()
				}
				imgui.EndMenu()
			}
			if imgui.BeginMenu(imguiw.T("Audio")) {
				if imgui.MenuItemBool(imguiw.T("AudioPlay")) {
					playAudio()
				}
				if imgui.MenuItemBool(imguiw.T("AudioPause")) {
					pauseAudio()
				}
				if imgui.MenuItemBool(imguiw.T("AudioStop")) {
					stopAudio()
				}
				if imgui.MenuItemBool(imguiw.RS("시작위치로 이동")) {
					goStartPosAudio()
				}
				imgui.EndMenu()
			}
			imgui.EndMenuBar()
		}
		widget.GetWaveformPlot().View()
	}
	imgui.End()
	imgui.PopStyleVar()
}

func openFile() {
	go func() {
		logger.Trace("[Event Callback] 파일 열기")

		filename, err := dialog.File().Filter("WAV Files", "wav").Load()
		if err != nil {
			logger.Errorf("파일 열기 실패 (filename=%v, err=%v)", filename, err)
		}
		logger.Infof("파일 경로 지정됨 (filename=%v)", filename)

		audio.Open(filename)
	}()
}

func playAudio() {
	go func() {
		logger.Trace("[Event Callback] 오디오 재생")
		audio.Play()
	}()
}

func pauseAudio() {
	go func() {
		logger.Trace("[Event Callback] 오디오 일시중지")
		audio.Pause()
	}()
}

func stopAudio() {
	go func() {
		logger.Trace("[Event Callback] 오디오 중지")
		audio.Stop()
	}()
}

func goStartPosAudio() {
	go func() {
		logger.Trace("[Event Callback] 오디오 초기 위치로 이동")
		audio.SetPosition(0)
	}()
}
