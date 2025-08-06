package indexador

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Indexador struct {
	DbConn        *pgx.Conn
	DbOpTimeout   time.Duration
	DbTxTimeout   time.Duration
	S3BucketName  string
	S3InitTimeout time.Duration
	S3OpTimeout   time.Duration
}

func (i *Indexador) ObtenerMaterias(ctx context.Context) ([]Materia, error) {
	var err error
	if err = i.initClienteS3(ctx); err != nil {
		return nil, err
	}

	var ofertas []OfertaMateriaSiu
	if ofertas, err = i.obtenerOfertasSiu(ctx); err != nil {
		return nil, err
	}

	var materias []Materia
	if materias, err = i.sincronizarConDb(ctx, ofertas); err != nil {
		return nil, err
	}

	return materias, nil
}
