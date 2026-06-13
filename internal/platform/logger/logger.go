package logger

import (
	"io"

	"github.com/username/project-name/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(appCfg config.AppConfig, obsCfg config.ObservabilityConfig) (*zap.Logger, io.Closer, error) {
	var zapCfg zap.Config
	if appCfg.Env == "development" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	baseLogger, err := zapCfg.Build()
	if err != nil {
		return nil, nil, err
	}

	lokiSyncer, err := newLokiWriteSyncer(appCfg, obsCfg.Loki)
	if err != nil {
		return nil, nil, err
	}

	if lokiSyncer == nil {
		return baseLogger, nil, nil
	}

	core := zapcore.NewTee(
		baseLogger.Core(),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapCfg.EncoderConfig),
			zapcore.AddSync(lokiSyncer),
			zapCfg.Level,
		),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapCfg.Level.Level()))

	return logger, lokiSyncer, nil
}
