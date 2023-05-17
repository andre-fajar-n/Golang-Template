package runtime

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func (r *Runtime) db() *Runtime {
	r.Logger.Info().Msg("Initiate connection to DB...")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		r.Cfg.DBHost,
		r.Cfg.DBUser,
		r.Cfg.DBPassword,
		r.Cfg.DBName,
		r.Cfg.DBPort,
	)

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	gormConfig := &gorm.Config{
		// enhance performance config
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 dbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error connect to DB")
		log.Panicf("Error connect to DB : %v", err)
	}

	r.Db = db

	r.Logger.Info().Msg("DB successfully connected")

	return r
}

func (r *Runtime) runMigration() {
	r.Logger.Info().Msg("Initiate db migration")

	r.Db.AutoMigrate()

	r.Logger.Info().Msg("Migrating db has been done")
}
