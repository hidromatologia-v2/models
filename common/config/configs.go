package config

import _ "github.com/sethvargo/go-envconfig"

type (
	Consumer struct {
		Station         string  `env:"STATION,required"`
		Consumer        string  `env:"CONSUMER,required"`
		Host            string  `env:"HOST,required"`
		Username        string  `env:"USERNAME,required"`
		Password        *string `env:"PASSWORD,noinit"`
		ConnectionToken *string `env:"CONN_TOKEN,noinit"`
	}
	Producer struct {
		Station         string  `env:"STATION,required"`
		Producer        string  `env:"PRODUCER,required"`
		Host            string  `env:"HOST,required"`
		Username        string  `env:"USERNAME,required"`
		Password        *string `env:"PASSWORD,noinit"`
		ConnectionToken *string `env:"CONN_TOKEN,noinit"`
	}
	Redis struct {
		Addr string `env:"ADDR,required"`
		DB   int    `env:"DB,required"`
	}
	SMTP struct {
		From     string  `env:"FROM,required"`
		Host     string  `env:"HOST,required"`
		Port     int     `env:"PORT,required"`
		Username *string `env:"USERNAME,noinit"`
		Password *string `env:"PASSWORD,noinit"`
		NoTLS    *bool   `env:"NO_TLS,noinit"`
	}
	Postgres struct {
		DSN string `env:"DSN,required"`
	}
	JWT struct {
		Secret string `env:"SECRET,required"`
	}
)
