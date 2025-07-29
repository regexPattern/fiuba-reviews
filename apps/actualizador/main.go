package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/patch"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/tui"
)

func main() {
	setupLogger()

	i := patch.Indexador{
		DbUrl:         os.Getenv("DATABASE_URL"),
		DbInitTimeout: time.Second * 3,
		DbOpTimeout:   time.Second * 10,
		S3BucketName:  os.Getenv("AWS_S3_BUCKET"),
		S3InitTimeout: time.Second * 3,
		S3OpTimeout:   time.Second * 10,
	}

	patches, err := i.GenerarPatches(context.Background())
	if err != nil {
		slog.Error("no se pudieron generar los patches de actualización")
		os.Exit(1)
	}

	for _, p := range patches[:10] {
		fmt.Println(p)
	}

	tui.ResolvePatches([]patch.Patch{})
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
