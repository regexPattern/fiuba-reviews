package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Db *pgxpool.Pool
	S3 *S3Config
}

func NewConfig() (Config, error) {
	var cfg Config

	db, err := newDbPool()
	if err != nil {
		return cfg, err
	}

	s3, err := newS3Client()
	if err != nil {
		return cfg, err
	}

	cfg = Config{
		Db: db,
		S3: s3,
	}

	return cfg, nil
}
