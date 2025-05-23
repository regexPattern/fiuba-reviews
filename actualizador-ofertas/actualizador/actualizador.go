package actualizador

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/config"
)

func GetPatches(cfg config.Config) ([]PatchMateriaOutput, error) {
	logger := log.Default()

	ofs, err := getOfertas(logger, cfg.S3)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ofertas de comisiones: %w", err)
	}

	ps, err := getPatchesMateriaOutput(logger, cfg.Db, ofs)
	if err != nil {
		return nil, fmt.Errorf(
			"error generando patches de actualización de materias: %w", err)
	}

	return ps, nil
}

func WritePatches(_ config.Config) error {
	return nil
}

func logErrRetMsg(logger *log.Logger, msg string, err error) error {
	logger.Helper()
	logger.Error(msg, "err", err)
	return errors.New(msg)
}
