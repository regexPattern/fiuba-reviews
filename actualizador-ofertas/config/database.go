package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbUrlEnvVarName string = "DATABASE_URL"

func newDbPool() (*pgxpool.Pool, error) {
	url, ok := os.LookupEnv(dbUrlEnvVarName)
	if !ok {
		return nil, fmt.Errorf("variable de entorno `%v` no configurada", dbUrlEnvVarName)
	}
	db, err := newDbPoolFromUrl(url)
	if err != nil {
		return nil, fmt.Errorf("error parseando string de conexión con la base de datos: %w", err)
	}
	if err := checkDbConn(db); err != nil {
		return nil, fmt.Errorf("error estableciendo conexión con la base de datos: %w", err)
	}
	return db, nil
}

func newDbPoolFromUrl(url string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	if pool, err := pgxpool.New(ctx, url); err != nil {
		return nil, err
	} else {
		return pool, nil
	}
}

func checkDbConn(db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		return err
	}
	return nil
}
