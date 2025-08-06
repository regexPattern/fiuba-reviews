package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/indexador"
	"github.com/regexPattern/fiuba-reviews/apps/actualizador/resolver"
)

func main() {
	setupLogger()

	conn, err := newDbConn(os.Getenv("DATABASE_URL"), time.Second*3)
	if err != nil {
		os.Exit(1)
	}

	i := indexador.Indexador{
		DbConn:        conn,
		DbOpTimeout:   time.Second * 10,
		DbTxTimeout:   time.Second * 10,
		S3BucketName:  os.Getenv("AWS_S3_BUCKET"),
		S3InitTimeout: time.Second * 3,
		S3OpTimeout:   time.Second * 10,
	}

	materias, err := i.ObtenerMaterias(context.Background())
	if err != nil {
		os.Exit(1)
	}

	if resolver.ResolverActualizaciones(conn, materias) != nil {
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
		TimeFormat:      time.TimeOnly,
		Level:           logLvl,
	})

	slog.SetDefault(slog.New(logger))
}

func newDbConn(dbUrl string, timeout time.Duration) (*pgx.Conn, error) {
	initctx, initcancel := context.WithTimeout(context.Background(), timeout)
	defer initcancel()

	conn, err := pgx.Connect(initctx, dbUrl)
	if err != nil {
		slog.Error(
			"error configurando pool de conexiones con la base de datos",
			"error",
			err,
		)
		return nil, err
	}

	slog.Info(
		"pool de conexiones con la base de datos configurado exitosamente",
	)

	return conn, nil
}
