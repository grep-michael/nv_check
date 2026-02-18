package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	gpubuilder "github.com/grep-michael/nv_check/Lib/GPUBuilder"
)

type OverviewRow struct {
	nameLabel  *widget.Label
	utilLabel  *widget.Label
	memLabel   *widget.Label
	tempLabel  *widget.Label
	powerLabel *widget.Label
}

type OverviewManager struct {
	vbox   *fyne.Container
	scroll *container.Scroll
	rows   []*OverviewRow
}

func BuildOverViewManager(gpus []gpubuilder.GPU) *OverviewManager {
	manager := &OverviewManager{}
	manager.buildVBox()
	manager.initRows(gpus)
	manager.buildScroll()
	return manager
}
func (manager *OverviewManager) buildScroll() {
	if manager.vbox == nil {
		manager.buildVBox()
	}
	manager.scroll = container.NewVScroll(manager.vbox)
}
func (manager *OverviewManager) buildVBox() {
	manager.vbox = container.NewVBox(manager.buildHeader(), widget.NewSeparator())
}
func (manager *OverviewManager) buildHeader() *fyne.Container {
	return container.NewGridWithColumns(5,
		widget.NewLabelWithStyle("GPU", fyne.TextAlignLeading, fyne.TextStyle{}),
		widget.NewLabelWithStyle("GPU Util", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Memory Used", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Temp", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Power Draw", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)
}
func (manager *OverviewManager) initRows(gpus []gpubuilder.GPU) {
	manager.rows = make([]*OverviewRow, len(gpus))
	if manager.vbox == nil {
		manager.buildVBox()
	}
	for i, gpu := range gpus {
		r := &OverviewRow{
			nameLabel:  widget.NewLabelWithStyle(gpuDisplayName(i, gpu), fyne.TextAlignCenter, fyne.TextStyle{}),
			utilLabel:  widget.NewLabelWithStyle(gpu.Utilization.UtilizationPercent, fyne.TextAlignCenter, fyne.TextStyle{}),
			memLabel:   widget.NewLabelWithStyle(memUsedDisplay(gpu), fyne.TextAlignCenter, fyne.TextStyle{}),
			tempLabel:  widget.NewLabelWithStyle(gpu.Temperature.Temp, fyne.TextAlignCenter, fyne.TextStyle{}),
			powerLabel: widget.NewLabelWithStyle(gpu.PowerReadings.InstantPowerDraw, fyne.TextAlignCenter, fyne.TextStyle{}),
		}
		manager.rows[i] = r

		row := container.NewGridWithColumns(5,
			r.nameLabel, r.utilLabel, r.memLabel, r.tempLabel, r.powerLabel,
		)
		manager.vbox.Add(row)
	}
}
func (manager *OverviewManager) Update(gpus []gpubuilder.GPU) {
	for i, gpu := range gpus {
		if i >= len(manager.rows) {
			break
		}
		r := manager.rows[i]
		r.nameLabel.SetText(gpuDisplayName(i, gpu))
		r.utilLabel.SetText(gpu.Utilization.UtilizationPercent)
		r.memLabel.SetText(memUsedDisplay(gpu))
		r.tempLabel.SetText(gpu.Temperature.Temp)
		r.powerLabel.SetText(gpu.PowerReadings.InstantPowerDraw)
	}
}
func (manager *OverviewManager) GetView() fyne.CanvasObject {
	return manager.scroll
}
