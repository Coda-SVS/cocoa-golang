package widget

import (
	"fmt"
	"sync"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/Kor-SVS/cocoa/src/audio"
	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/ui/imguiw"
	"github.com/Kor-SVS/cocoa/src/util"
	"github.com/sasha-s/go-deadlock"
	"github.com/zyedidia/generic/list"
)

var (
	plotOnce   sync.Once
	plotLogger *log.Logger
)

type SubPlot struct {
	dataPlot imguiw.PlotWidget
	ratio    float32
}

func NewSubPlot(dataPlot imguiw.PlotWidget) *SubPlot {
	return &SubPlot{
		dataPlot: dataPlot,
		ratio:    1.0,
	}
}

type Plot struct {
	logger    *log.Logger
	dataPlots *list.List[*SubPlot]
	wg        *sync.WaitGroup
	mtx       *deadlock.Mutex

	axisXLimitMax float64

	isFitAudioStreamPos bool

	audioStreamPosID           int32
	audioStreamPos             float64
	audioStreamPos_out_clicked bool
	audioStreamPos_out_hovered bool
	audioStreamPos_held        bool
	audioStreamPos_IsPaused    bool

	// temp state
	isFitRequest         bool
	row_ratios           []float32
	plotDrawEndEventArgs *util.PlotDrawEndEventArgs
}

func NewPlot() *Plot {
	plotOnce.Do(func() {
		logOption := log.NewLoggerOption()
		logOption.Prefix = "[plot]"
		plotLogger = logger.NewSimpleLogger(logOption)

		plotLogger.Trace("Plot init...")
	})

	p := &Plot{}
	p.logger = plotLogger
	p.dataPlots = list.New[*SubPlot]()
	p.wg = &sync.WaitGroup{}
	p.mtx = &deadlock.Mutex{}
	p.row_ratios = make([]float32, 0)

	p.eventHandler_AudioStreamChanged()
	p.axisXLimitMax = DefaultAxisXLimitMax

	p.isFitAudioStreamPos = true

	p.plotDrawEndEventArgs = &util.PlotDrawEndEventArgs{}

	return p
}

func (p *Plot) EditDataPlotList(fn func(*list.List[*SubPlot])) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	fn(p.dataPlots)
}

func (p *Plot) View() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	firstSubPlotNode := p.dataPlots.Front
	if firstSubPlotNode == nil {
		return
	}

	// 서브 플롯 데이터 준비
	p.row_ratios = p.row_ratios[0:0]
	subPlotCount := 0
	for dpn := firstSubPlotNode; dpn != nil; dpn = dpn.Next {
		subPlot := dpn.Value
		subPlotCount++
		p.row_ratios = append(p.row_ratios, subPlot.ratio)
		go subPlot.dataPlot.UpdateData()
	}

	if imgui.PlotBeginSubplotsV(
		imguiw.RS("SubPlot##widget"),
		int32(subPlotCount),
		1,
		imgui.Vec2{X: -1, Y: -1},
		imgui.PlotSubplotFlagsLinkCols|
			imgui.PlotSubplotFlagsNoLegend|
			imgui.PlotSubplotFlagsNoTitle|
			imgui.PlotSubplotFlagsNoMenus,
		&p.row_ratios,
		nil,
	) {
		var plotDrawEndEventArgs *util.PlotDrawEndEventArgs
		isFitRequest := p.isFitRequest
		isFitAudioStreamPos := p.isFitAudioStreamPos

		col := -1
		for dpn := firstSubPlotNode; dpn != nil; dpn = dpn.Next {
			col++
			subPlot := dpn.Value
			subPlot.ratio = p.row_ratios[col]
			dataPlot := subPlot.dataPlot

			if imgui.PlotBeginPlotV(
				imguiw.RS(fmt.Sprintf("Plot##%v", col)),
				imgui.Vec2{X: -1, Y: -1},
				imgui.PlotFlagsNoLegend|
					imgui.PlotFlagsNoTitle|
					imgui.PlotFlagsNoMouseText|
					imgui.PlotFlagsNoMenus,
			) {
				// Plot Setup
				isLastSubPlot := col == subPlotCount-1
				dataPlot.PlotSetup(imguiw.PlotSetupArgs{
					AxisXLimitMax:       p.axisXLimitMax,
					IsLastSubPlot:       isLastSubPlot,
					IsFitRequest:        isFitRequest,
					IsFitAudioStreamPos: isFitAudioStreamPos,
				})
				imgui.PlotSetupFinish()

				// Draw Plot
				dataPlot.Plot()

				// Draw Stream Pos Line
				p.audioStreamPosDragLine()

				if imgui.PlotIsPlotHovered() && imgui.IsMouseDoubleClicked(imgui.MouseButtonLeft) {
					p.isFitRequest = true
				}

				if plotDrawEndEventArgs == nil {
					plotSize := imgui.PlotGetPlotSize()
					plotPos := imgui.PlotGetPlotPos()

					plotEndSize := plotSize.Add(plotPos)

					plotPointStart := imgui.PlotPixelsToPlotFloatV(plotPos.X, plotPos.Y, imgui.AxisX1, imgui.AxisY1).X
					plotPointEnd := imgui.PlotPixelsToPlotFloatV(plotEndSize.X, plotEndSize.Y, imgui.AxisX1, imgui.AxisY1).X

					plotDrawEndEventArgs = p.plotDrawEndEventArgs
					plotDrawEndEventArgs.PlotPixelXStart = float64(plotPos.X)
					plotDrawEndEventArgs.PlotPixelXEnd = float64(plotEndSize.X)
					plotDrawEndEventArgs.PlotPointStart = plotPointStart
					plotDrawEndEventArgs.PlotPointEnd = plotPointEnd
				}

				imgui.PlotEndPlot()
			}
		}

		if isFitRequest {
			p.isFitRequest = false
		}

		if plotDrawEndEventArgs != nil {
			for dpn := firstSubPlotNode; dpn != nil; dpn = dpn.Next {
				p.wg.Add(1)
				go func(pw imguiw.PlotWidget) {
					defer p.wg.Done()
					pw.EventHandler(*plotDrawEndEventArgs)
				}(dpn.Value.dataPlot)
			}
			p.wg.Wait()
		}
		imgui.PlotEndSubplots()
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
					p.axisXLimitMax = DefaultAxisXLimitMax
				}
				p.isFitRequest = true
			case audio.EnumAudioStreamClosed:
				p.axisXLimitMax = DefaultAxisXLimitMax
				imgui.PlotBustPlotCache()
			}
		}
	}()
}
