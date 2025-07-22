package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patcher"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/resolver"
)

func main() {
	setupLogger()

	genCtx, cancelGenCtx := context.WithTimeout(context.Background(), time.Millisecond*3000)
	defer cancelGenCtx()

	proposed, err := patcher.GeneratePatches(genCtx)
	if err != nil {
		slog.Error("no se pudieron generar los patches de actualización")
		os.Exit(1)
	}

	resolved := resolver.ResolvePatches(proposed)

	applyCtx, cancelApplyCtx := context.WithTimeout(context.Background(), time.Millisecond*3000)
	defer cancelApplyCtx()

	if patcher.ApplyPatches(applyCtx, resolved) != nil {
		slog.Error("no se pudieron aplicar los patches de actualización")
		os.Exit(1)
	}
}

func setupLogger() {
	logLvl := log.InfoLevel
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logLvl = log.DebugLevel
	}

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Level:           logLvl,
	})

	slog.SetDefault(slog.New(logger))
}
