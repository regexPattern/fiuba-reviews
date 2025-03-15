package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Db *pgxpool.Pool
	S3 *S3Config
}

func NewConfig() (*Config, error) {
	db, err := newDbPool()
	if err != nil {
		return nil, err
	}

	s3, err := newS3Client()
	if err != nil {
		return nil, err
	}

	c := &Config{
		Db: db,
		S3: s3,
	}
	return c, nil
}
