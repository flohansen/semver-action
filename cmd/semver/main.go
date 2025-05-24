package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flohansen/semver/internal/app"
)

func main() {
	var cfg app.Config
	flag.StringVar(&cfg.Token, "token", "", "The token used for authentication")
	flag.Parse()

	action := app.NewAction(cfg)
	if err := action.Run(app.SignalContext()); err != nil {
		fmt.Printf("error running action: %s", err)
		os.Exit(1)
	}
}
