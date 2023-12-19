package imguiw

type Widget interface {
	View()
}

type Window interface {
	Widget

	Title() string

	IsOpen() bool
	SetIsOpen(bool)
}

type PlotWidget interface {
	Title() string
	Plot()
	UpdateData()
	EventHandler(eventArgs any)
	IsDisposed() bool
}
