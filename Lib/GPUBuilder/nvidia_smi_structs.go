package gpubuilder

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Root struct {
	XMLName xml.Name `xml:"nvidia_smi_log"` // match your actual root tag
	GPUs    []GPU    `xml:"gpu"`
}

type GPU struct {
	XMLName            xml.Name               `xml:"gpu"`
	ID                 string                 `xml:"id,attr"`
	UUID               string                 `xml:"uuid"`
	Name               string                 `xml:"product_name"`
	Brand              string                 `xml:"product_brand"`
	PartNumb           string                 `xml:"gpu_part_number"`
	PerformanceState   string                 `xml:"performance_state"`
	FanSpeed           string                 `xml:"fan_speed"`
	Utilization        UtilizationBlock       `xml:"utilization"`
	Memory             MemoryBlock            `xml:"fb_memory_usage"`
	Temperature        TemperatureBlock       `xml:"temperature"`
	SupportTargetTemp  SupportTargetTempBlock `xml:"supported_gpu_target_temp"`
	PowerReadings      PowerReadingsBlock     `xml:"gpu_power_readings"`
	CurrentClockSpeeds ClockSpeedBlock        `xml:"clocks"`
	MaxClockSpeeds     ClockSpeedBlock        `xml:"max_clocks"`
}

type ClockSpeedBlock struct {
	GraphicsClock string `xml:"graphics_clock"`
	SMClock       string `xml:"sm_clock"`
	MemoryClock   string `xml:"mem_clock"`
	VideoClock    string `xml:"video_clock"`
}

type PowerReadingsBlock struct {
	PowerState       string `xml:"power_state"`
	AveragePowerDraw string `xml:"average_power_draw"`
	InstantPowerDraw string `xml:"instant_power_draw"`
	PowerLimit       string `xml:"current_power_limit"`
}

type SupportTargetTempBlock struct {
	MinTarget string `xml:"gpu_target_temp_min"`
	Maxtarget string `xml:"gpu_target_temp_max"`
}

type TemperatureBlock struct {
	Temp              string `xml:"gpu_temp"`
	TLimit            string `xml:"gpu_temp_tlimit"`
	MaxThreshold      string `xml:"gpu_temp_max_threshold"`
	SlowThresHold     string `xml:"gpu_temp_slow_threshold"`
	MaxGPUThresHold   string `xml:"gpu_temp_max_gpu_threshold"`
	TargetTemperature string `xml:"gpu_target_temperature"`
	MemoryTemp        string `xml:"memory_temp"`
}

type MemoryBlock struct {
	Total    MemoryNumber `xml:"total"`
	Reserved MemoryNumber `xml:"reserved"`
	Used     MemoryNumber `xml:"used"`
	Free     MemoryNumber `xml:"free"`
}
type MemoryNumber int

func (mem *MemoryNumber) MarshalJSON() ([]byte, error) {
	dict := []string{"Bytes", "KiB", "MiB"} //, "GiB"} //, "PiB", "man what are you doing with this much ram"}
	temp := float64(int(*mem))
	index := 0
	for ; temp > 1024 && index < len(dict)-1; index++ {
		temp = temp / 1024
	}
	s := fmt.Sprintf("%0.2f %s", temp, dict[index])
	return json.Marshal(s)
}

func (mem *MemoryNumber) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	dict := map[string]int{
		"KiB": 1024,
		"MiB": 1024 * 1024,
		"GiB": 1000 * 1024,
	}
	var size_raw string
	if err := d.DecodeElement(&size_raw, &start); err != nil {
		return err
	}

	size_list := strings.Split(size_raw, " ")
	if len(size_list) != 2 {
		return fmt.Errorf("Split Size into more/less than 2 elements: %+v", size_list)
	}
	size_s := size_list[0]
	unit_s := size_list[1]

	size_i, err := strconv.Atoi(size_s)
	if err != nil {
		return err
	}

	multiplyer, ok := dict[unit_s]
	if !ok {
		log.Printf("Memory unit not in dictionary")
		multiplyer = 1
	}
	size_i = size_i * multiplyer
	*mem = MemoryNumber(size_i)
	return nil
}

type UtilizationBlock struct {
	UtilizationPercent       string `xml:"gpu_util"`
	MemoryUtilizarionPercent string `xml:"memory_util"`
}
