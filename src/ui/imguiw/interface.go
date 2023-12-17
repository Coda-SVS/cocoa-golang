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
