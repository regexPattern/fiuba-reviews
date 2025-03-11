package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dbUrlEnvVarName string = "DATABASE_URL"

var db *pgxpool.Pool

// conectarDb establece la conexión con la base de datos.
func conectarDb(logger *log.Logger) error {
	dbUrlEnv, ok := os.LookupEnv(dbUrlEnvVarName)
	if !ok {
		return fmt.Errorf("variable de entorno `%v` no configurada", dbUrlEnvVarName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var err error
	db, err = pgxpool.New(ctx, dbUrlEnv)

	if err != nil {
		logger.Error(err)
		return errors.New("error configurando conexión con la base de datos")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if db.Ping(ctx) != nil {
		logger.Error(err)
		return errors.New("error estableciendo conexión con la base de datos")
	}

	logger.Info("conexión con la base de datos establecida exitosamente")

	return nil
}
