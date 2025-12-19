package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func main() {
	logger := log.New(os.Stderr)

	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logger.SetLevel(log.DebugLevel)
	}

	slog.SetDefault(slog.New(logger))

	dbUrl := os.Getenv("DATABASE_URL")
	host := os.Getenv("BACKEND_HOST")
	port := os.Getenv("BACKEND_PORT")

	addr := net.JoinHostPort(host, port)

	if err := run(dbUrl, addr); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(dbUrl, addr string) error {
	conn, err := pgx.Connect(context.TODO(), dbUrl)
	if err != nil {
		return fmt.Errorf("error estableciendo conexión con la base de datos: %w", err)
	}

	slog.Info("conexión establecida con la base de datos")

	patches, err := genPatchesMaterias(conn)
	if err != nil {
		return fmt.Errorf("error generando patches de materias: %w", err)
	}

	if err := startServer(conn, addr, patches); err != nil {
		return fmt.Errorf(
			"error iniciando servidor de patches de materias: %w",
			err,
		)
	}

	return nil
}

func genPatchesMaterias(conn *pgx.Conn) (map[string]patchMateria, error) {
	ofertas, err := getOfertasMaterias(conn)
	if err != nil {
		return nil, fmt.Errorf(
			"error obteniendo ofertas de comisiones de materias: %w",
			err,
		)
	}

	codMaterias := make([]string, 0, len(ofertas))
	nombresMaterias := make([]string, 0, len(ofertas))

	for cod, om := range ofertas {
		codMaterias = append(codMaterias, cod)
		nombresMaterias = append(nombresMaterias, om.Nombre)
	}

	if err := syncDb(conn, codMaterias, nombresMaterias); err != nil {
		return nil, fmt.Errorf(
			"error sincronizando materias de la base de datos con el siu: %w",
			err,
		)
	}

	patches, err := buildPatchesMaterias(conn, codMaterias, ofertas)
	if err != nil {
		return nil, fmt.Errorf(
			"error construyendo patches de actualización de materias: %w",
			err,
		)
	}

	return patches, nil
}

func startServer(conn *pgx.Conn, addr string, patches map[string]patchMateria) error {
	http.HandleFunc("GET /patches", func(w http.ResponseWriter, _ *http.Request) {
		handleGetAllPatches(w, patches)
	})
	http.HandleFunc("GET /patches/{codigoMateria}", func(w http.ResponseWriter, r *http.Request) {
		handleGetPatchMateria(w, r, conn, patches)
	})
	http.HandleFunc("PATCH /patches/{codigoMateria}", func(w http.ResponseWriter, r *http.Request) {
		handleAplicarPatchMateria(w, r, conn, patches)
	})

	slog.Info(fmt.Sprintf("servidor escuchando peticiones en dirección %v", addr))

	return http.ListenAndServe(addr, nil)
}
