package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDbPool(logger *log.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var err error
	db, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		logger.Error(err)
		return errors.New("error estableciendo conexión con la base de datos")
	}

	logger.Debug("establecida conexión con la base de datos")

	return nil
}
