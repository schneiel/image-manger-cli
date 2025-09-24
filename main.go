// Package main provides the entry point for the Image Manager CLI application.
package main

import (
	"embed"
	stdlog "log"
	"os"

	cli "github.com/schneiel/ImageManagerGo/internal"
	"github.com/schneiel/ImageManagerGo/internal/di"
)

//go:embed locales
var localesFS embed.FS

func main() {
	container, err := di.NewContainer(os.Args, &localesFS)
	if err != nil {
		stdlog.Fatal(err)
	}

	// Create command executor with injected dependencies from DI container using builder.
	executor, err := cli.NewCommandExecutorBuilder().
		WithArgs(os.Args).
		WithLocalizer(container.Localizer).
		WithFileReader(container.FileReader).
		WithParser(container.Parser).
		WithConfig(container.Config).
		WithHandlers(container.SortHandler, container.DedupHandler).
		WithFlagSetups(container.SortFlagSetup, container.DedupFlagSetup).
		WithCommandLocalizer(container.CommandLocalizer).
		Build()
	if err != nil {
		stdlog.Fatal(err)
	}

	err = executor.Execute()
	if err != nil {
		stdlog.Fatal(err)
	}
}
