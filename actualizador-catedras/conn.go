package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func initDbConn() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var err error
	conn, err = pgx.Connect(ctx, os.Getenv("DATABASE_URL"))

	return err
}
