package gpubuilder

import (
	"encoding/xml"
	"log"
	"os/exec"
)

func BuildGPUS() (gpuList []GPU, err error) {
	xml_bytes, err := exec.Command("nvidia-smi", "-q", "-x").Output()
	if err != nil {
		log.Printf("Error Executing nvidia-smi -q -x: %+v\n", err)
		return
	}

	var root Root
	err = xml.Unmarshal(xml_bytes, &root)
	if err != nil {
		log.Printf("Error Unmarshaling: %v\n", err)
		return
	}
	gpuList = root.GPUs
	return
}
