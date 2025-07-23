package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/resolver"
)

func main() {
	setupLogger()

	g := patch.Generador{
		DbUrl:         os.Getenv("DATABASE_URL"),
		DbTimeout:     time.Second * 3,
		S3BucketName:  os.Getenv("AWS_S3_BUCKET"),
		S3InitTimeout: time.Second * 3,
		S3Timeout:     time.Second * 3,
	}

	proposed, err := g.GenerarPatches(context.Background())
	if err != nil {
		slog.Error("no se pudieron generar los patches de actualización")
		os.Exit(1)
	}

	if proposed == nil {
		return
	}

	resolved := resolver.ResolvePatches(proposed)

	if patch.ApplyPatches(context.Background(), resolved) != nil {
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
