package imguiw

type Widget interface {
	View()
	Close()
}

type Window interface {
	Widget

	Title() string

	IsOpen() bool
	SetIsOpen(bool)
}
