package imguiw

import (
	"sync"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/sasha-s/go-deadlock"
)

var Context *ImguiWContext

type ImguiWContext struct {
	imDPI     *ImguiDPI
	imBackend imgui.Backend[imgui.GLFWWindowFlags]
	FontAtlas *FontAtlas
	context   *imgui.Context
	mutex     *deadlock.RWMutex
	waitGroup *sync.WaitGroup
	idCounter int
}

func (ic *ImguiWContext) Backend() imgui.Backend[imgui.GLFWWindowFlags] {
	return ic.imBackend
}

func (ic *ImguiWContext) ID() int {
	ic.idCounter++
	return ic.idCounter
}

func (ic *ImguiWContext) IO() *imgui.IO {
	return imgui.CurrentIO()
}

func (ic *ImguiWContext) WaitGroup() *sync.WaitGroup {
	return ic.waitGroup
}

func (ic *ImguiWContext) Mutex() *deadlock.RWMutex {
	return ic.mutex
}
