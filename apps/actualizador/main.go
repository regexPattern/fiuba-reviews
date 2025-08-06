package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/resolver"
)

func main() {
	setupLogger()

	i := indexador.Indexador{
		DbUrl:         os.Getenv("DATABASE_URL"),
		DbInitTimeout: time.Second * 3,
		DbOpTimeout:   time.Second * 10,
		S3BucketName:  os.Getenv("AWS_S3_BUCKET"),
		S3InitTimeout: time.Second * 3,
		S3OpTimeout:   time.Second * 10,
	}

	ofertas, err := i.ObtenerMaterias(context.Background())
	if err != nil {
		os.Exit(1)
	}

	if len(ofertas) == 0 {
		slog.Info("no hay materias por actualizar")
		return
	}

	resolver.ResolverPatches(nil)
}

func setupLogger() {
	logLvl := log.InfoLevel
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logLvl = log.DebugLevel
	}

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           logLvl,
	})

	slog.SetDefault(slog.New(logger))
}
