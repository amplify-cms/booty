package main

import (
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/dep/orchestrator"
)

var (
	// best variable name
	conductor dep.Commander
	version   = ""
	revision  = ""
)

func init() {
	conductor = orchestrator.NewOrchestrator("booty", version, revision)
}

func main() {
	logger := conductor.Logger()
	if err := conductor.Command().Execute(); err != nil {
		logger.Errorf("error: %v", err)
	}
}
