package actualizador

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/regexPattern/fiuba-reviews/actualizador-ofertas/config"
)

func GetPatches(cfg *config.Config) ([]PatchMateriaOutput, error) {
	logger := log.Default()

	_, err := getOfertas(logger, cfg.S3)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ofertas de comisiones: %w", err)
	}

	return nil, nil
}

func WritePatches(_ *config.Config) error {
	return nil
}

func wrapErrorMsg(logger *log.Logger, msg string, err error) error {
	logger.Error(msg, "err", err)
	return errors.New(msg)
}
