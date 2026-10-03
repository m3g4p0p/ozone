package main

import (
	"context"
	"flag"
	"log"

	"m3g4p0p/ozone/examples/util"
	"m3g4p0p/ozone/ozone"
	"m3g4p0p/ozone/ozone/think"
)

func run() error {
	agent := &ozone.Agent{
		Name:  "main",
		Model: "qwen3.5:2b",
		Think: think.False,
	}

	return util.RunLoop(context.Background(), agent)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
