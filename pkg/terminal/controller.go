package terminal

import (
	"fmt"
	"runtime"

	"github.com/Zioyi/temi/pkg"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	terminalWidth     = 120
	heapAllocBarCount = 6
)

type controller struct {
	Grid *ui.Grid

	HeapObjectsSparkline  *widgets.Sparkline
	HeapObjectsSparkGroup *widgets.SparklineGroup
	HeapObjectsData       *pkg.StatRing

	SysText       *widgets.Paragraph
	GCCPUFraction *widgets.Gauge

	HeapAllocBarChar     *widgets.BarChart
	HeapAllocBarCharData *pkg.StatRing

	HeapPie *widgets.PieChart
}

func (p *controller) Resize() {
	p.resize()
	ui.Render(p.Grid)
}

func (p *controller) resize() {
	_, h := ui.TerminalDimensions()
	p.Grid.SetRect(0, 0, terminalWidth, h)
}

func (p *controller) Render(data *runtime.MemStats) {
	p.HeapObjectsData.Push(data.HeapObjects)
	p.HeapObjectsSparkline.Data = p.HeapObjectsData.NormalizedData()
	p.HeapObjectsSparkGroup.Title = fmt.Sprintf("HeapObjects, live heap object count: %d", data.HeapObjects)

	p.SysText.Text = fmt.Sprintf(byteCountBinary(data.Sys))

	fNormalize := func() int {
		f := data.GCCPUFraction
		if f < 0.01 && f > 0 {
			for f < 1 {
				f = f * 10.0
			}
		}
		return int(f)
	}
	p.GCCPUFraction.Percent = fNormalize()
	p.GCCPUFraction.Label = fmt.Sprintf("%.2f%%", data.GCCPUFraction*100)

	p.HeapAllocBarCharData.Push(data.HeapAlloc)
	p.HeapAllocBarChar.Data = p.HeapAllocBarCharData.Data()
	p.HeapAllocBarChar.Labels = nil
	for _, v := range p.HeapAllocBarChar.Data {
		p.HeapAllocBarChar.Labels = append(p.HeapAllocBarChar.Labels, byteCountBinary(uint64(v)))
	}

	p.HeapPie.Data = []float64{float64(data.HeapIdle), float64(data.HeapInuse)}

	ui.Render(p.Grid)
}

func (p *controller) initUI() {
	p.resize()

	p.HeapObjectsSparkline.LineColor = ui.Color(89) // DeepPink4
	p.HeapObjectsSparkGroup = widgets.NewSparklineGroup(p.HeapObjectsSparkline)

	p.SysText.Title = "Sys, the total bytes of allocated heap objects"
	p.SysText.PaddingLeft = 25
	p.SysText.PaddingTop = 1

	p.HeapAllocBarChar.BarGap = 2
	p.HeapAllocBarChar.BarWidth = 8
	p.HeapAllocBarChar.Title = "HeapAlloc, bytes of allocated heap objects"
	p.HeapAllocBarChar.NumFormatter = func(f float64) string { return "" }

	p.GCCPUFraction.Title = "GCCPUFraction 0%~100%"
	p.GCCPUFraction.BarColor = ui.Color(50) // Cyan2

	p.HeapPie.Title = "HeapInuse vs HeadIdle"
	p.HeapPie.LabelFormatter = func(idx int, _ float64) string { return []string{"Idle", "Inuse"}[idx] }

	p.Grid.Set(
		ui.NewRow(.2, p.HeapObjectsSparkGroup),
		ui.NewRow(.8,
			ui.NewCol(.5,
				ui.NewRow(.2, p.SysText),
				ui.NewRow(.2, p.GCCPUFraction),
				ui.NewRow(.6, p.HeapAllocBarChar),
			),
			ui.NewCol(.5, p.HeapPie),
		),
	)
}

func newController() *controller {
	ctl := &controller{
		Grid: ui.NewGrid(),

		HeapObjectsSparkline: widgets.NewSparkline(),
		HeapObjectsData:      pkg.NewChartRing(terminalWidth),

		SysText:       widgets.NewParagraph(),
		GCCPUFraction: widgets.NewGauge(),

		HeapAllocBarChar:     widgets.NewBarChart(),
		HeapAllocBarCharData: pkg.NewChartRing(heapAllocBarCount),

		HeapPie: widgets.NewPieChart(),
	}

	ctl.initUI()

	return ctl
}
