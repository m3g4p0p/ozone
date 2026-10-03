package main

import (
	"context"
	"flag"
	"log"
	"os"

	"m3g4p0p/ozone/examples/util"
	"m3g4p0p/ozone/ozone"
	"m3g4p0p/ozone/ozone/think"
	"m3g4p0p/ozone/ozone/tools"
)

func run() error {
	listDir := tools.MustNewTool(
		"list_dir",
		"list files in a directory",
		func(ctx context.Context, input struct {
			Dirname string `json:"dirname" jsonschema:"name of the directory to read"`
		},
		) ([]string, error) {
			entries, err := os.ReadDir(input.Dirname)
			if err != nil {
				return nil, tools.LLMError(err)
			}
			var names []string
			for _, e := range entries {
				names = append(names, e.Name())
			}
			return names, nil
		},
	)

	readFile := tools.MustNewTool(
		"read_file",
		"read the contents of a file",
		func(ctx context.Context, input struct {
			Filename string `json:"filename" jsonschema:"path to the file to read"`
		},
		) (string, error) {
			data, err := os.ReadFile(input.Filename)
			if err != nil {
				return "", tools.LLMError(err)
			}
			return string(data), nil
		},
	)

	agent := &ozone.Agent{
		Name:  "main",
		Model: "qwen3.5:2b",
		Think: think.False,
		Toolset: &tools.Deferred{
			Name:        "activate_file_tools",
			Description: "enabled access to list_dir and read_file",
			Provider:    tools.Toolsets{listDir, readFile},
			Enabled:     tools.HasActivateMessage,
		},
	}

	return util.RunLoop(context.Background(), agent)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
