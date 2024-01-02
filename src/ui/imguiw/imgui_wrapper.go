package imguiw

import (
	"sync"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/util"
)

var (
	logger *log.Logger

	beforeDestroyContextCallback func()
)

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[imguiw]"
	logger = log.RootLogger().NewSimpleLogger(logOption)

	logger.Trace("Imguiw init...")
}

func InitImgui(title string, width, height int) {
	Context = &ImguiWContext{}

	Context.waitGroup = util.GetWaitGroup()

	Context.context = imgui.CreateContext()
	imgui.SetCurrentContext(Context.context)
	Context.imBackend = imgui.CreateBackend(imgui.NewGLFWBackend())
	Context.imDPI = NewImguiDPI(nil, Context.context, nil)
	Context.FontAtlas = newFontAtlas()

	io := imgui.CurrentIO()

	io.SetConfigFlags(imgui.ConfigFlagsDpiEnableScaleViewports)
	io.SetBackendFlags(imgui.BackendFlagsRendererHasVtxOffset)

	io.SetIniFilename("")

	Context.imBackend.SetTargetFPS(60)
	Context.imBackend.SetBeforeRenderHook(beforeRender)
	Context.imBackend.SetAfterCreateContextHook(afterCreateContext)
	Context.imBackend.SetBeforeDestroyContextHook(beforeDestroyContext)

	// Create font
	fonts := Context.IO().Fonts()
	fonts.Clear()
	if len(Context.FontAtlas.defaultFonts) == 0 {
		fonts.AddFontDefault()
		fontTextureImg, w, h, _ := fonts.GetTextureDataAsRGBA32()
		tex := Context.imBackend.CreateTexture(fontTextureImg, int(w), int(h))
		fonts.SetTexID(tex)
		fonts.SetTexReady(true)
	} else {
		Context.FontAtlas.shouldRebuildFontAtlas = true
	}

	Context.imBackend.CreateWindow(title, width, height)

	logger.Trace("(Call) InitImgui")
}

func Run(fn func()) {
	Context.imBackend.Run(fn)
}

func beforeRender() {
	Context.FontAtlas.rebuildFontAtlas()
}

func afterCreateContext() {
	imgui.PlotCreateContext()
	imgui.ImNodesCreateContext()

	logger.Trace("(Call) afterCreateContext")
}

func SetBeforeDestroyContextCallback(f func()) {
	beforeDestroyContextCallback = f
}

func beforeDestroyContext() {
	if beforeDestroyContextCallback != nil {
		Context.waitGroup.Add(1)
		go func() {
			beforeDestroyContextCallback()
			Context.waitGroup.Done()
		}()
	}

	if waitTimeout(Context.waitGroup, time.Duration(time.Second*5)) {
		logger.Error("서브루틴 종료실패 (time-out)")
	}

	imgui.ImNodesDestroyContext()
	imgui.PlotDestroyContext()

	logger.Trace("(Call) beforeDestroyContext")
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
