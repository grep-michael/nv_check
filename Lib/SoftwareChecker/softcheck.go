package softwarechecker

import (
	"os/exec"
)

type Result struct {
	Name      string
	Installed bool
	Path      string
}

func Check(name string) Result {
	path, err := exec.LookPath(name)
	if err != nil {
		return Result{Name: name, Installed: false}
	}
	return Result{Name: name, Installed: true, Path: path}
}
func CheckAll(names []string) map[string]Result {
	results := make(map[string]Result, len(names))
	for _, name := range names {
		results[name] = Check(name)
	}
	return results
}
func AllInstalled(names []string) bool {
	for _, name := range names {
		if !Check(name).Installed {
			return false
		}
	}
	return true
}
