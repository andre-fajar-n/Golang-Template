package go_template

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-openapi/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewRuntime() *Runtime {
	rt := new(Runtime)

	rt = rt.logger()

	rt = rt.config()

	rt = rt.db()

	rt.runMigration()

	return rt
}

type Runtime struct {
	Db     *gorm.DB
	Cfg    config
	Logger zerolog.Logger
}

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

func (r *Runtime) config() *Runtime {
	r.Logger.Info().Msg("Initiate read env...")

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	for _, v := range os.Environ() {
		temp := strings.Split(v, "=")
		viper.BindEnv(temp[0])
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {

		default:
			r.Logger.Error().Err(err).Msg("Failed load config")
			log.Fatalf("Failed load config : %v", err)

		case viper.ConfigFileNotFoundError:
			err = viper.SafeWriteConfig()
			if err != nil {
				switch err.(type) {

				case viper.ConfigFileAlreadyExistsError:

				default:
					r.Logger.Error().Err(err).Msg("Failed SafeWriteConfig config")
					log.Fatalf("Failed SafeWriteConfig config : %v", err)
				}
			}
		}
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		r.Logger.Error().Err(err).Msg("Failed unmarshal config")
		log.Fatalf("Failed unmarshal config : %v", err)
	}

	r.Cfg = cfg

	r.Logger.Info().Msg("Success reading env")

	return r
}

func (r *Runtime) logger() *Runtime {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}).With().Timestamp().Caller().Logger()

	r.Logger = logger

	return r
}

func (r *Runtime) runMigration() {
	r.Logger.Info().Msg("Initiate db migration")

	r.Db.AutoMigrate()

	r.Logger.Info().Msg("Migrating db has been done")
}

func (r *Runtime) SetError(code int, msg string, args ...interface{}) error {
	return errors.New(int32(code), msg, args...)
}

func (r *Runtime) GetError(err error) errors.Error {
	if v, ok := err.(errors.Error); ok {
		return v
	}

	return errors.New(http.StatusInternalServerError, err.Error())
}
