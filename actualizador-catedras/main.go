package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := initLogger()
	slog.SetDefault(logger)

	materias, err := getMateriasOfertasComisiones()
	if err != nil {
		slog.Error("Error obteniendo las ofertas de comisiones disponibles", "error", err)
	}

	fmt.Println(materias)

	getCodigosMaterias()
}

func initLogger() *slog.Logger {
	var level slog.Level

	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
