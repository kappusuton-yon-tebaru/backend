package logger

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New(cfg *config.Config) (*Logger, error) {
	var logger *zap.Logger
	var err error

	if cfg.Development {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return &Logger{
		logger,
	}, nil
}
