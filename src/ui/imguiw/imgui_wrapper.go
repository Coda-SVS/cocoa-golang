package imguiw

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/log"
)

var logger *log.Logger

func init() {
	logOption := log.NewLoggerOption()
	logOption.Prefix = "[imguiw]"
	logWriter := log.NewLogWriter(nil, nil, nil, nil)
	logger = log.NewLogger(logOption, logWriter)

	logger.Trace("Imguiw init...")
}

func InitImgui() {
	Context = &ImguiWContext{}
	Context.context = imgui.CreateContext()
	Context.imBackend = imgui.CreateBackend(imgui.NewGLFWBackend())
	Context.imDPI = NewImguiDPI(nil, Context.context, nil)
	Context.FontAtlas = newFontAtlas()

	io := imgui.CurrentIO()

	io.SetConfigFlags(imgui.BackendFlagsRendererHasVtxOffset)
	io.SetBackendFlags(imgui.BackendFlagsRendererHasVtxOffset)

	io.SetIniFilename("")

	// logger.Tracef("imgui.CurrentIO().IniFilename()=%v", imgui.CurrentIO().IniFilename())
	// logger.Tracef("io.IniFilename()=%v", io.IniFilename())

	Context.imBackend.SetTargetFPS(60)
	Context.imBackend.SetAfterCreateContextHook(afterCreateContext)
	Context.imBackend.SetBeforeDestroyContextHook(beforeDestroyContext)
	Context.imBackend.SetBeforeRenderHook(beforeRender)

	// Context.imBackend.SetWindowFlags(imgui.GLFWWindowFlagsVisible, 0)
	// Context.imBackend.SetWindowFlags(imgui.GLFWWindowFlagsFloating, 0)
	// Context.imBackend.SetWindowFlags(imgui.GLFWWindowFlagsTransparent, 0)
	// io.SetConfigViewportsNoAutoMerge(true)

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

	logger.Trace("(Call) InitImgui")
}

func afterCreateContext() {
	imgui.PlotCreateContext()

	logger.Trace("(Call) afterCreateContext")
}

func beforeDestroyContext() {
	imgui.PlotDestroyContext()

	logger.Trace("(Call) beforeDestroyContext")
}

func beforeRender() {
	Context.FontAtlas.rebuildFontAtlas()
}
