package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func main() {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	slog.SetDefault(slog.New(logger))

	dbUrl := os.Getenv("DATABASE_URL")
	host := os.Getenv("BACKEND_HOST")
	port := os.Getenv("BACKEND_PORT")

	addr := net.JoinHostPort(host, port)

	if err := run(dbUrl, addr); err != nil {
		slog.Error("aplicacion_fallida", "error", err)
		os.Exit(1)
	}
}

func run(dbUrl, addr string) error {
	conn, err := pgx.Connect(context.TODO(), dbUrl)
	if err != nil {
		return fmt.Errorf("error estableciendo conexión con la base de datos: %w", err)
	}

	slog.Info("conexion_con_db_establecida")

	patches, err := getPatchesMaterias(conn)
	if err != nil {
		return fmt.Errorf("error generando patches de materias: %w", err)
	}

	if err := iniciarServidor(conn, addr, patches); err != nil {
		return fmt.Errorf(
			"error iniciando servidor de patches de materias: %w",
			err,
		)
	}

	return nil
}
