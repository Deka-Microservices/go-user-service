package database

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type deka_db_logger struct {
	level logger.LogLevel
}

func newLogger() *deka_db_logger {
	return &deka_db_logger{
		level: logger.Info,
	}
}

func (l *deka_db_logger) LogMode(level logger.LogLevel) logger.Interface {
	new_logger := *l
	new_logger.level = level
	return &new_logger
}

func (l deka_db_logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		log.Info().Any("data", data).Msg(msg)
	}
}

func (l deka_db_logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		log.Warn().Any("data", data).Msg(msg)
	}
}

func (l deka_db_logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		log.Error().Any("data", data).Msg(msg)
	}
}

func (l deka_db_logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	if l.level <= logger.Silent {
		return
	}

	switch {
	case err != nil && l.level >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		if rows == -1 {
			log.Error().Str("filename", utils.FileWithLineNum()).Err(err).Str("sql", sql).Msg("database trace")
		} else {
			log.Error().Str("filename", utils.FileWithLineNum()).Err(err).Str("sql", sql).Int64("rows", rows).Msg("database trace")
		}
	case l.level == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Debug().Str("filename", utils.FileWithLineNum()).Str("sql", sql).Msg("database trace")
		} else {
			log.Debug().Str("filename", utils.FileWithLineNum()).Str("sql", sql).Int64("rows", rows).Msg("database trace")
		}
	}

}
