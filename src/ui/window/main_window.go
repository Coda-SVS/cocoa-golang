package window

import (
	g "github.com/AllenDang/giu"
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
		),

		g.SplitLayout(g.DirectionVertical, &State.LeftSidePanelPos,
			// 좌측 패널
			g.Layout{
				g.Label("좌측 패널"),
			},
			// 우측 패널
			g.Layout{
				g.Label("우측 패널"),
			},
		),
	)
}

func openFile() {
	go func() {
		logger.Trace("[Event Callback] 파일 열기")

		filename, err := dialog.File().Filter("WAV Files", "wav").Load()
		logger.Errorf("filename=%v, err=%v", filename, err)
	}()
}
