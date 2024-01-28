package main

import (
	"fmt"
	"twoBinaries"

	"github.com/ppacer/core/exec"
)

func main() {
	const port = 8080
	dags := twoBinaries.SetupDags()
	executor := exec.New(fmt.Sprintf("http://localhost:%d", port), nil)
	executor.Start(dags)
}
