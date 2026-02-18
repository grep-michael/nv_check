package gpubuilder

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestIntMinBasic(t *testing.T) {
	gpus, err := BuildGPUS()
	if err != nil {
		log.Println(err)
	}
	printobj(gpus)
}

func printobj(obj any) {
	js, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(js))
}
