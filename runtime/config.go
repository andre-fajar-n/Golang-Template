package runtime

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type (
	config struct {
		DBHost     string `mapstructure:"DB_HOST"`
		DBUser     string `mapstructure:"DB_USER"`
		DBPassword string `mapstructure:"DB_PASSWORD"`
		DBName     string `mapstructure:"DB_NAME"`
		DBPort     int    `mapstructure:"DB_PORT"`

		JwtSecret string `mapstructure:"JWT_SECRET"`
		JwtExp    int    `mapstructure:"JWT_EXP"`
	}
)

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
