package database

import (
	"github.com/deka-microservices/go-user-service/internal/consts"
	"github.com/deka-microservices/go-user-service/pkg/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"

	"gorm.io/gorm"
)

const SQLITE_FALLBACK_FILE = "sqlite.db"

func NewUserDB() *gorm.DB {
	db_string := viper.GetString(consts.CONFIG_DSN)

	var gormDB *gorm.DB
	config := &gorm.Config{
		Logger: newLogger(),
	}

	var err error

	if len(db_string) == 0 {
		log.Warn().Msg("DSN string is empty. Switching to SQLite3")
		gormDB, err = gorm.Open(sqlite.Open(db_string), config)
		if err != nil {
			log.Fatal().Err(err).Msg("database connection fail")

		}
	} else {
		gormDB, err = gorm.Open(postgres.Open(db_string), config)
		if err != nil {
			log.Fatal().Err(err).Msg("database connection fail")
		}
	}

	return gormDB
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(models.User{})
}

type UserQuerier interface {
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (gen.T, error)

	// SELECT * FROM @@table WHERE username=@username
	GetByUsername(username string) (gen.T, error)
}
