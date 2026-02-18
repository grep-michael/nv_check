package gui

//AI Generated, good enough for now

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	gpubuilder "github.com/grep-michael/nv_check/Lib/GPUBuilder"
)

// RunApp starts the Fyne GPU monitor application.
func RunAppAi() {
	a := app.New()
	w := a.NewWindow("GPU Monitor")
	//w.Resize(fyne.NewSize(1000, 700))

	// Initial load
	gpus, err := gpubuilder.BuildGPUS()
	if err != nil {
		gpus = []gpubuilder.GPU{}
	}

	state := &AppState{
		gpus: gpus,
	}

	tabs := buildTabs(state)
	w.SetContent(tabs)

	// Background refresh goroutine
	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			fresh, err := gpubuilder.BuildGPUS()
			if err != nil {
				continue
			}
			state.mu.Lock()
			state.gpus = fresh
			state.mu.Unlock()

			// Refresh all bound labels on the UI thread
			fyne.Do(func() {
				state.mu.Lock()
				defer state.mu.Unlock()
				for i, gpu := range state.gpus {
					if i < len(state.bindings) {
						state.bindings[i].update(gpu)
					}
				}
				// Update overview cards too
				state.overviewRefresh()
			})
		}
	}()

	w.ShowAndRun()
}

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

type AppState struct {
	mu              sync.Mutex
	gpus            []gpubuilder.GPU
	bindings        []*gpuBindings
	overviewRefresh func()
}

// ---------------------------------------------------------------------------
// Tab builder
// ---------------------------------------------------------------------------

func buildTabs(state *AppState) *container.AppTabs {
	tabs := container.NewAppTabs()

	// Overview tab
	overviewContent, refreshFn := buildOverviewTab(state)
	state.overviewRefresh = refreshFn
	tabs.Append(container.NewTabItem("Overview", overviewContent))

	// Per-GPU detail tabs
	state.bindings = make([]*gpuBindings, len(state.gpus))
	for i, gpu := range state.gpus {
		b := newGPUBindings(gpu)
		state.bindings[i] = b
		label := fmt.Sprintf("GPU %d", i)
		if gpu.Name != "" {
			label = gpu.Name
		}
		tabs.Append(container.NewTabItem(label, buildDetailTab(b)))
	}

	return tabs
}

// ---------------------------------------------------------------------------
// Overview tab
// ---------------------------------------------------------------------------

func buildOverviewTab(state *AppState) (fyne.CanvasObject, func()) {
	// We store widget.Label pointers per GPU so we can refresh them.
	type overviewRow struct {
		nameLabel  *widget.Label
		utilLabel  *widget.Label
		memLabel   *widget.Label
		tempLabel  *widget.Label
		powerLabel *widget.Label
	}

	rows := make([]*overviewRow, len(state.gpus))

	header := container.NewGridWithColumns(5,
		widget.NewLabelWithStyle("GPU", fyne.TextAlignLeading, fyne.TextStyle{}),
		widget.NewLabelWithStyle("GPU Util", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Memory Used", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Temp", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Power Draw", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	grid := container.NewVBox(header, widget.NewSeparator())

	for i, gpu := range state.gpus {
		r := &overviewRow{
			nameLabel:  widget.NewLabelWithStyle(gpuDisplayName(i, gpu), fyne.TextAlignCenter, fyne.TextStyle{}),
			utilLabel:  widget.NewLabelWithStyle(gpu.Utilization.UtilizationPercent, fyne.TextAlignCenter, fyne.TextStyle{}),
			memLabel:   widget.NewLabelWithStyle(memUsedDisplay(gpu), fyne.TextAlignCenter, fyne.TextStyle{}),
			tempLabel:  widget.NewLabelWithStyle(gpu.Temperature.Temp, fyne.TextAlignCenter, fyne.TextStyle{}),
			powerLabel: widget.NewLabelWithStyle(gpu.PowerReadings.InstantPowerDraw, fyne.TextAlignCenter, fyne.TextStyle{}),
		}
		rows[i] = r

		row := container.NewGridWithColumns(5,
			r.nameLabel, r.utilLabel, r.memLabel, r.tempLabel, r.powerLabel,
		)
		grid.Add(row)
	}

	scroll := container.NewVScroll(grid)

	refresh := func() {
		for i, gpu := range state.gpus {
			if i >= len(rows) {
				break
			}
			r := rows[i]
			r.nameLabel.SetText(gpuDisplayName(i, gpu))
			r.utilLabel.SetText(gpu.Utilization.UtilizationPercent)
			r.memLabel.SetText(memUsedDisplay(gpu))
			r.tempLabel.SetText(gpu.Temperature.Temp)
			r.powerLabel.SetText(gpu.PowerReadings.InstantPowerDraw)
		}
	}

	return scroll, refresh
}

// ---------------------------------------------------------------------------
// Detail tab
// ---------------------------------------------------------------------------

// gpuBindings holds all labels for a single GPU detail view so they can be
// updated without rebuilding the widget tree.
type gpuBindings struct {
	// Identity
	name, brand, UUID, partNum, id *widget.Label
	// Performance
	perfState, fanSpeed *widget.Label
	// Utilization
	gpuUtil, memUtil *widget.Label
	// Memory
	memTotal, memUsed, memFree, memReserved *widget.Label
	// Temperature
	temp, tLimit, maxThresh, slowThresh, targetTemp, memTemp *widget.Label
	// Power
	powerState, avgPower, instPower, powerLimit *widget.Label
	// Clocks (current)
	curGfx, curSM, curMem, curVid *widget.Label
	// Clocks (max)
	maxGfx, maxSM, maxMem, maxVid *widget.Label
}

func newGPUBindings(gpu gpubuilder.GPU) *gpuBindings {
	b := &gpuBindings{}
	b.populate(gpu)
	return b
}

func (b *gpuBindings) populate(gpu gpubuilder.GPU) {
	lbl := func(s string) *widget.Label { return widget.NewLabel(s) }
	b.name = lbl(gpu.Name)
	b.brand = lbl(gpu.Brand)
	b.UUID = lbl(gpu.UUID)
	b.partNum = lbl(gpu.PartNumb)
	b.id = lbl(gpu.ID)
	b.perfState = lbl(gpu.PerformanceState)
	b.fanSpeed = lbl(gpu.FanSpeed)
	b.gpuUtil = lbl(gpu.Utilization.UtilizationPercent)
	b.memUtil = lbl(gpu.Utilization.MemoryUtilizarionPercent)
	b.memTotal = lbl(memStr(gpu.Memory.Total))
	b.memUsed = lbl(memStr(gpu.Memory.Used))
	b.memFree = lbl(memStr(gpu.Memory.Free))
	b.memReserved = lbl(memStr(gpu.Memory.Reserved))
	b.temp = lbl(gpu.Temperature.Temp)
	b.tLimit = lbl(gpu.Temperature.TLimit)
	b.maxThresh = lbl(gpu.Temperature.MaxThreshold)
	b.slowThresh = lbl(gpu.Temperature.SlowThresHold)
	b.targetTemp = lbl(gpu.Temperature.TargetTemperature)
	b.memTemp = lbl(gpu.Temperature.MemoryTemp)
	b.powerState = lbl(gpu.PowerReadings.PowerState)
	b.avgPower = lbl(gpu.PowerReadings.AveragePowerDraw)
	b.instPower = lbl(gpu.PowerReadings.InstantPowerDraw)
	b.powerLimit = lbl(gpu.PowerReadings.PowerLimit)
	b.curGfx = lbl(gpu.CurrentClockSpeeds.GraphicsClock)
	b.curSM = lbl(gpu.CurrentClockSpeeds.SMClock)
	b.curMem = lbl(gpu.CurrentClockSpeeds.MemoryClock)
	b.curVid = lbl(gpu.CurrentClockSpeeds.VideoClock)
	b.maxGfx = lbl(gpu.MaxClockSpeeds.GraphicsClock)
	b.maxSM = lbl(gpu.MaxClockSpeeds.SMClock)
	b.maxMem = lbl(gpu.MaxClockSpeeds.MemoryClock)
	b.maxVid = lbl(gpu.MaxClockSpeeds.VideoClock)
}

func (b *gpuBindings) update(gpu gpubuilder.GPU) {
	set := func(l *widget.Label, s string) { l.SetText(s) }
	set(b.name, gpu.Name)
	set(b.brand, gpu.Brand)
	set(b.UUID, gpu.UUID)
	set(b.partNum, gpu.PartNumb)
	set(b.id, gpu.ID)
	set(b.perfState, gpu.PerformanceState)
	set(b.fanSpeed, gpu.FanSpeed)
	set(b.gpuUtil, gpu.Utilization.UtilizationPercent)
	set(b.memUtil, gpu.Utilization.MemoryUtilizarionPercent)
	set(b.memTotal, memStr(gpu.Memory.Total))
	set(b.memUsed, memStr(gpu.Memory.Used))
	set(b.memFree, memStr(gpu.Memory.Free))
	set(b.memReserved, memStr(gpu.Memory.Reserved))
	set(b.temp, gpu.Temperature.Temp)
	set(b.tLimit, gpu.Temperature.TLimit)
	set(b.maxThresh, gpu.Temperature.MaxThreshold)
	set(b.slowThresh, gpu.Temperature.SlowThresHold)
	set(b.targetTemp, gpu.Temperature.TargetTemperature)
	set(b.memTemp, gpu.Temperature.MemoryTemp)
	set(b.powerState, gpu.PowerReadings.PowerState)
	set(b.avgPower, gpu.PowerReadings.AveragePowerDraw)
	set(b.instPower, gpu.PowerReadings.InstantPowerDraw)
	set(b.powerLimit, gpu.PowerReadings.PowerLimit)
	set(b.curGfx, gpu.CurrentClockSpeeds.GraphicsClock)
	set(b.curSM, gpu.CurrentClockSpeeds.SMClock)
	set(b.curMem, gpu.CurrentClockSpeeds.MemoryClock)
	set(b.curVid, gpu.CurrentClockSpeeds.VideoClock)
	set(b.maxGfx, gpu.MaxClockSpeeds.GraphicsClock)
	set(b.maxSM, gpu.MaxClockSpeeds.SMClock)
	set(b.maxMem, gpu.MaxClockSpeeds.MemoryClock)
	set(b.maxVid, gpu.MaxClockSpeeds.VideoClock)
}

func buildDetailTab(b *gpuBindings) fyne.CanvasObject {
	row := func(key string, val *widget.Label) *fyne.Container {
		return container.NewGridWithColumns(2,
			widget.NewLabelWithStyle(key, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			val,
		)
	}

	identity := widget.NewCard("Identity", "",
		container.NewVBox(
			row("Name:", b.name),
			row("Brand:", b.brand),
			row("UUID:", b.UUID),
			row("Part Number:", b.partNum),
			row("ID:", b.id),
			row("Performance State:", b.perfState),
			row("Fan Speed:", b.fanSpeed),
		),
	)

	utilization := widget.NewCard("Utilization", "",
		container.NewVBox(
			row("GPU Util:", b.gpuUtil),
			row("Memory Util:", b.memUtil),
		),
	)

	memory := widget.NewCard("Memory", "",
		container.NewVBox(
			row("Total:", b.memTotal),
			row("Used:", b.memUsed),
			row("Free:", b.memFree),
			row("Reserved:", b.memReserved),
		),
	)

	temperature := widget.NewCard("Temperature", "",
		container.NewVBox(
			row("GPU Temp:", b.temp),
			row("T-Limit:", b.tLimit),
			row("Max Threshold:", b.maxThresh),
			row("Slow Threshold:", b.slowThresh),
			row("Target Temp:", b.targetTemp),
			row("Memory Temp:", b.memTemp),
		),
	)

	power := widget.NewCard("Power", "",
		container.NewVBox(
			row("Power State:", b.powerState),
			row("Instant Draw:", b.instPower),
			row("Average Draw:", b.avgPower),
			row("Power Limit:", b.powerLimit),
		),
	)

	clocks := widget.NewCard("Clock Speeds", "",
		container.NewVBox(
			widget.NewLabelWithStyle("Current", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true}),
			row("  Graphics:", b.curGfx),
			row("  SM:", b.curSM),
			row("  Memory:", b.curMem),
			row("  Video:", b.curVid),
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Max", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true}),
			row("  Graphics:", b.maxGfx),
			row("  SM:", b.maxSM),
			row("  Memory:", b.maxMem),
			row("  Video:", b.maxVid),
		),
	)

	content := container.NewVBox(identity, utilization, memory, temperature, power, clocks)
	return container.NewVScroll(content)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func gpuDisplayName(i int, gpu gpubuilder.GPU) string {
	if gpu.Name != "" {
		return fmt.Sprintf("[%d] %s", i, gpu.Name)
	}
	return fmt.Sprintf("GPU %d", i)
}

func memUsedDisplay(gpu gpubuilder.GPU) string {
	return fmt.Sprintf("%s / %s", memStr(gpu.Memory.Used), memStr(gpu.Memory.Total))
}

// memStr renders a MemoryNumber using its own JSON marshalling logic so we
// get a human-readable string without needing to duplicate the conversion.
func memStr(m gpubuilder.MemoryNumber) string {
	dict := []string{"Bytes", "KiB", "MiB"}
	temp := float64(int(m))
	index := 0
	for ; temp > 1024 && index < len(dict)-1; index++ {
		temp /= 1024
	}
	return fmt.Sprintf("%0.2f %s", temp, dict[index])
}
