package window

import (
	g "github.com/AllenDang/giu"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/ui/plot"
	"github.com/sqweek/dialog"
)

type MainWindowState struct {
	LeftSidePanelPos float32
}

func newMainWindowState() (state MainWindowState) {
	state = MainWindowState{
		LeftSidePanelPos: 400,
	}
	return state
}

var (
	State MainWindowState = newMainWindowState()
)

func MainWindowGUILoop() {
	windowMutex.RLock()
	defer windowMutex.RUnlock()

	g.SingleWindowWithMenuBar().Layout(
		// 메뉴바
		g.MenuBar().Layout(
			g.Menu("파일").Layout(
				g.MenuItem("열기").OnClick(openFile),
				// g.MenuItem("Save"),
				// // You could add any kind of widget here, not just menu item.
				// g.Menu("Save as ...").Layout(
				// 	g.MenuItem("Excel file"),
				// 	g.MenuItem("CSV file"),
				// 	g.Button("Button inside menu"),
				// ),
			),
			g.Menu("오디오").Layout(
				g.MenuItem("재생").OnClick(playAudio),
				g.MenuItem("일시중지").OnClick(pauseAudio),
				g.MenuItem("중지").OnClick(stopAudio),
				g.MenuItem("처음 위치로").OnClick(goStartPosAudio),
			),
		),

		g.SplitLayout(g.DirectionVertical, &State.LeftSidePanelPos,
			// 좌측 패널
			g.Layout{
				g.Label("좌측 패널"),
			},
			// 우측 패널
			g.Layout{
				g.Label("우측 패널"),
				plot.WaveformGUILoop(),
			},
		),
	)
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

		audio.Play()
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
