package go_template

type (
	config struct {
		DBHost     string `mapstructure:"DB_HOST"`
		DBUser     string `mapstructure:"DB_USER"`
		DBPassword string `mapstructure:"DB_PASSWORD"`
		DBName     string `mapstructure:"DB_NAME"`
		DBPort     int    `mapstructure:"DB_PORT"`

		JwtSecret string `mapstructure:"JWT_SECRET"`
	}
)
