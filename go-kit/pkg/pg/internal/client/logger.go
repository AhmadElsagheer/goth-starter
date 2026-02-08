package client

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

type logger struct {
	logger *zap.SugaredLogger
}

func (l *logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	log := l.logger
	keyAndValues := make([]interface{}, 0, len(data)*2)
	for k, v := range data {
		keyAndValues = append(keyAndValues, k, v)
	}

	switch level {
	case tracelog.LogLevelTrace:
		log.Debugw(msg, append(keyAndValues, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case tracelog.LogLevelDebug:
		log.Debugw(msg, keyAndValues...)
	case tracelog.LogLevelInfo:
		log.Infow(msg, keyAndValues...)
	case tracelog.LogLevelWarn:
		log.Warnw(msg, keyAndValues...)
	case tracelog.LogLevelError:
		log.Errorw(msg, keyAndValues...)
	default:
		log.Errorw(msg, append(keyAndValues, zap.Stringer("PGX_LOG_LEVEL", level))...)
	}
}
