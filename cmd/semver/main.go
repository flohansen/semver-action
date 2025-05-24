package main

import (
	"fmt"
	"os"

	"github.com/flohansen/semver/internal/app"
)

func main() {
	action := app.NewAction()
	if err := action.Run(app.SignalContext()); err != nil {
		fmt.Printf("error running action: %s", err)
		os.Exit(1)
	}
}
