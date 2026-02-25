package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	gpubuilder "github.com/grep-michael/nv_check/Lib/GPUBuilder"
	"os/exec"
	"strings"
)

type PCITab struct {
	rows          []string
	lspci         []string
	filtered      []string
	list          *widget.List
	Status        string
	CurrentFilter string
}

func BuildPCITab() *PCITab {
	manager := &PCITab{}
	manager.initPCIDevices()
	return manager
}

func (tab *PCITab) Update([]gpubuilder.GPU) {
}

func (tab *PCITab) GetView() fyne.CanvasObject {
	header := tab.buildHeader()
	tab.list = tab.buildList()
	tab.refreshBody()
	return container.NewBorder(header, nil, nil, nil, tab.list)
}

func (tab *PCITab) buildList() *widget.List {
	return widget.NewList(
		func() int {
			return len(tab.filtered)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(tab.filtered[i])
		},
	)
}

func (tab *PCITab) refreshBody() {
	tab.filtered = nil

	for _, line := range tab.lspci {
		if tab.CurrentFilter == "" {
			tab.filtered = append(tab.filtered, line)
			continue
		}
		if strings.Contains(
			strings.ToLower(line),
			strings.ToLower(tab.CurrentFilter),
		) {
			tab.filtered = append(tab.filtered, line)
		}
	}

	tab.list.Refresh()
}

func (tab *PCITab) buildHeader() *fyne.Container {
	btn1 := widget.NewButton("Ethernet", func() {
		tab.CurrentFilter = "ethernet"
		tab.refreshBody()
	})
	btn2 := widget.NewButton("VGA", func() {
		tab.CurrentFilter = "vga"
		tab.refreshBody()
	})
	btn3 := widget.NewButton("clear", func() {
		tab.CurrentFilter = ""
		tab.refreshBody()
	})

	textInput := widget.NewEntry()
	textInput.SetPlaceHolder("Custom search")
	textInput.OnChanged = func(s string) {
		tab.CurrentFilter = s
		tab.refreshBody()
	}

	return container.NewBorder(
		nil, nil,
		container.NewHBox(btn1, btn2, btn3),
		nil,
		textInput,
	)
}

func (tab *PCITab) initPCIDevices() {
	stdout, err := exec.Command("lspci").Output()
	if err != nil {
		tab.Status = fmt.Sprintf("Failed to get lspci: %+v", err)
		return
	}

	tab.lspci = nil
	for line := range strings.Lines(string(stdout)) {
		tab.lspci = append(tab.lspci, line)
	}
}
