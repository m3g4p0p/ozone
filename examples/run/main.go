package main

import (
	"context"
	"flag"
	"log"

	"m3g4p0p/ozone/examples/util"
	"m3g4p0p/ozone/ozone"
	"m3g4p0p/ozone/ozone/tools"
)

type GetWeatherParams struct {
	Location string `json:"location" jsonschema:"the location to get the weather for"`
}

func run() error {
	tool := tools.MustNewTool(
		"get_weather",
		"get the weather for the provided location",
		func(_ context.Context, params GetWeatherParams) (string, error) {
			return "mostly sunny", nil
		},
	)

	agent := &ozone.Agent{
		Name:    "main",
		Model:   "qwen3.5:2b",
		Toolset: tools.Toolsets{tool},
	}

	return util.RunSingle(
		context.Background(),
		agent,
		flag.Arg(0),
	)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
