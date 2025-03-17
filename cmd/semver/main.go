package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flohansen/semver/internal/app"
)

func main() {
	flags := app.ActionFlags{}
	flag.StringVar(&flags.OutputName, "output-name", "new-version", "The name of the Actions output containing the new version")
	flag.Parse()

	action := app.NewAction(flags)
	if err := action.Run(app.SignalContext()); err != nil {
		fmt.Printf("error running action: %s", err)
		os.Exit(1)
	}
}
