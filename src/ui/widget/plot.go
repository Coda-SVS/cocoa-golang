package widget

import (
	"sync"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
	"github.com/Kor-SVS/cocoa/src/util"
)

var (
	plotOnce     sync.Once
	plotInstance *Plot
	plotLogger   *log.Logger
)

type Plot struct {
	logger      *log.Logger
	plotWidgets map[string]*imguiw.PlotWidget
	wg          *sync.WaitGroup

	isFitRequest  bool
	axisXLimitMax float64

	audioStreamPosID           int32
	audioStreamPos             float64
	audioStreamPos_out_clicked bool
	audioStreamPos_out_hovered bool
	audioStreamPos_held        bool
	audioStreamPos_IsPaused    bool
}

func NewPlot() *Plot {
	plotOnce.Do(func() {
		logOption := log.NewLoggerOption()
		logOption.Prefix = "[plot]"
		plotLogger = logger.NewSimpleLogger(logOption)

		plotLogger.Trace("Plot init...")
	})

	plotInstance = &Plot{}
	plotInstance.logger = plotLogger
	plotInstance.plotWidgets = make(map[string]*imguiw.PlotWidget)
	plotInstance.wg = &sync.WaitGroup{}

	plotInstance.eventHandler_AudioStreamChanged()
	plotInstance.axisXLimitMax = DefaultAxisXLimitMax

	return plotInstance
}

func (p *Plot) AddPlot(plotWidget imguiw.PlotWidget) {
	if plotWidget == nil {
		return
	}

	p.plotWidgets[plotWidget.Title()] = &plotWidget
}

func (p *Plot) RemoveDisposedPlotData() {
	if len(p.plotWidgets) == 0 {
		return
	}

	for titleKey, plotDataValue := range p.plotWidgets {
		if (*plotDataValue).IsDisposed() {
			delete(p.plotWidgets, titleKey)
		}
	}
}

func (p *Plot) View() {
	p.RemoveDisposedPlotData()

	if len(p.plotWidgets) == 0 {
		return
	}

	if imgui.PlotBeginPlotV(
		imguiw.T("Plot"),
		imgui.Vec2{X: -1, Y: -1},
		imgui.PlotFlagsNoLegend|imgui.PlotFlagsNoTitle|imgui.PlotFlagsNoMenus|imgui.PlotFlagsNoMouseText,
	) {
		for _, dataPlot := range p.plotWidgets {
			p.wg.Add(1)
			go func(pw *imguiw.PlotWidget) {
				defer p.wg.Done()
				(*pw).UpdateData()
			}(dataPlot)
		}
		p.wg.Wait()

		imgui.PlotSetupAxisV(
			imgui.AxisX1,
			"PlotX",
			imgui.PlotAxisFlags(imgui.PlotAxisFlagsNoLabel),
		)
		imgui.PlotSetupAxisV(
			imgui.AxisY1,
			"PlotY",
			imgui.PlotAxisFlags(imgui.PlotAxisFlagsLock|imgui.PlotAxisFlagsNoTickLabels|imgui.PlotAxisFlagsNoGridLines|imgui.PlotAxisFlagsNoLabel),
		)

		imgui.PlotSetupAxisLimitsConstraints(imgui.AxisX1, 0, p.axisXLimitMax)
		imgui.PlotSetupAxisLimitsV(imgui.AxisY1, -0.5, 0.5, imgui.PlotCondAlways)

		if p.isFitRequest {
			imgui.PlotSetupAxisLimitsV(imgui.AxisX1, 0, p.axisXLimitMax, imgui.PlotCondAlways)
			p.isFitRequest = false
		}

		imgui.PlotSetupLock()

		for _, dataPlot := range p.plotWidgets {
			(*dataPlot).Plot()
		}

		p.audioStreamPosDragLine()

		if imgui.PlotIsPlotHovered() && imgui.IsMouseDoubleClicked(imgui.MouseButtonLeft) {
			p.isFitRequest = true
		}

		plotSize := imgui.PlotGetPlotSize()
		plotPos := imgui.PlotGetPlotPos()

		plotEndSize := plotSize.Add(plotPos)

		plotPointStart := imgui.PlotPixelsToPlotFloatV(plotPos.X, plotPos.Y, imgui.AxisX1, imgui.AxisY1).X
		plotPointEnd := imgui.PlotPixelsToPlotFloatV(plotEndSize.X, plotEndSize.Y, imgui.AxisX1, imgui.AxisY1).X
		plotDrawEndEventArgs := util.PlotDrawEndEventArgs{
			PlotPointStart: plotPointStart,
			PlotPointEnd:   plotPointEnd,
		}

		for _, dataPlot := range p.plotWidgets {
			p.wg.Add(1)
			go func(pw *imguiw.PlotWidget) {
				defer p.wg.Done()
				(*pw).EventHandler(plotDrawEndEventArgs)
			}(dataPlot)
		}
		p.wg.Wait()

		imgui.PlotEndPlot()
	}
}

// 재생 위치 표시 및 조작
func (p *Plot) audioStreamPosDragLine() {
	if !audio.IsAudioLoaded() {
		return
	}

	p.audioStreamPos = audio.Position().Seconds()
	_audioStreamPos := p.audioStreamPos
	imgui.PlotDragLineXV(
		p.audioStreamPosID,
		&p.audioStreamPos,
		imgui.NewVec4(255, 0, 0, 255),
		1,
		imgui.PlotDragToolFlagsDelayed,
		&p.audioStreamPos_out_clicked,
		&p.audioStreamPos_out_hovered,
		&p.audioStreamPos_held,
	)

	if p.audioStreamPos_held {
		if !p.audioStreamPos_IsPaused && audio.IsRunning() {
			audio.Pause()
			p.audioStreamPos_IsPaused = true
		}
		if _audioStreamPos != p.audioStreamPos {
			audio.SetPosition(time.Duration(p.audioStreamPos * float64(time.Second)))
		}
	}

	if p.audioStreamPos_out_clicked && p.audioStreamPos_IsPaused {
		audio.Play()
		p.audioStreamPos_IsPaused = false
	}
}

// 오디오 스트림의 이벤트 수신
func (p *Plot) eventHandler_AudioStreamChanged() {
	go func() {
		msgChan := audio.AudioStreamBroker().Subscribe()

		for msg := range msgChan {
			p.logger.Tracef("[Callback] AudioStreamChanged (msg=%v)", msg)
			switch msg {
			case audio.EnumAudioStreamOpen:
				p.axisXLimitMax = audio.Duration().Seconds()
				if p.axisXLimitMax == 0 {
					plotInstance.axisXLimitMax = DefaultAxisXLimitMax
				}
				p.isFitRequest = true
			case audio.EnumAudioStreamClosed:
				plotInstance.axisXLimitMax = DefaultAxisXLimitMax
				imgui.PlotBustPlotCache()
			}
		}
	}()
}
