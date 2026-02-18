package softwarechecker

import (
	"log"
)

func CheckRequiredSoftware() bool {
	required := []string{"nvidia-smi", "nvcc", "hipcc", "gcc", "apt"}
	for _, cmd := range required {
		result := Check(cmd)
		if !result.Installed {
			log.Printf("Missing package that provides %s\n", result.Name)
			return false
		}
	}
	return true
}
