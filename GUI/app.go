package gui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	gpubuilder "github.com/grep-michael/nv_check/Lib/GPUBuilder"
)

func RunApp() {
	a := app.New()
	w := a.NewWindow("GPU Monitor")
	w.Resize(fyne.NewSize(1000, 700))

	nv_app := NewApp()

	overviewManager := BuildOverViewManager(nv_app.gpus)
	pci_tab := BuildPCITab()
	nv_app.RegisterTab("Overview", overviewManager)
	nv_app.RegisterTab("PCI Devices", pci_tab)

	w.SetContent(nv_app.tabs)
	go nv_app.UpdateLoop()
	w.ShowAndRun()
}

type UpdateableTable interface {
	Update([]gpubuilder.GPU)
	GetView() fyne.CanvasObject
}

type App struct {
	mu          sync.Mutex
	gpus        []gpubuilder.GPU
	tabs        *container.AppTabs
	tabMap      map[string]UpdateableTable
	SelectedTab string
	bindings    []*gpuBindings
}

func NewApp() *App {
	app := &App{
		tabMap: make(map[string]UpdateableTable),
	}

	gpus, err := gpubuilder.BuildGPUS()
	if err != nil {
		gpus = []gpubuilder.GPU{}
	}
	app.gpus = gpus
	app.tabs = container.NewAppTabs()
	app.tabs.OnSelected = func(ti *container.TabItem) {
		fmt.Printf("Selected tab: %s\n", ti.Text)
		app.SelectedTab = ti.Text
	}

	return app
}

func (app *App) GetSelectedTab() string {
	selected := app.tabs.Selected()
	if selected == nil {
		return ""
	}
	return selected.Text
}
func (app *App) UpdateLoop() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		fresh, err := gpubuilder.BuildGPUS()
		if err != nil {
			continue
		}
		app.mu.Lock()
		app.gpus = fresh
		tmp := app.gpus
		app.mu.Unlock()

		fyne.Do(func() {
			tabName := app.GetSelectedTab()
			tab, ok := app.tabMap[tabName]
			if ok {
				tab.Update(tmp)
			}
		})
	}
}
func (app *App) RegisterTab(text string, tabItem UpdateableTable) {
	app.tabMap[text] = tabItem
	app.tabs.Append(container.NewTabItem(text, tabItem.GetView()))
}
