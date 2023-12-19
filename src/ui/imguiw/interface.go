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
	Plot()
	UpdateData()
	EventHandler(eventArgs any)
	IsDisposed() bool
}
