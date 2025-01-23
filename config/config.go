package config

import (
	"fmt"

	"github.com/kochabonline/kit/config"
	"github.com/kochabonline/kit/log"
)

var Cfg = new(Config)

func init() {
	c := config.NewConfig(config.Option{Target: Cfg})
	if err := c.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := c.WatchConfig(); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Http  Http
	Mysql Mysql
	Redis Redis
	Mongo Mongo
	Kafka Kafka
	Jwt   Jwt
	Email Email
	Audit Audit
}

type Http struct {
	Host    string `default:"localhost"`
	Port    int    `default:"8080"`
	Level   string `default:"release"`
	Swagger struct {
		Enabled bool
	}
	Health struct {
		Enabled bool
	}
	Metrics struct {
		Enabled bool
	}
}

type Mysql struct {
	Host     string `default:"localhost"`
	Port     int    `default:"3306"`
	User     string `default:"root"`
	Password string
	Database string
	Level    string `default:"silent"`
}

type Redis struct {
	Host     string `default:"localhost"`
	Port     int    `default:"6379"`
	Password string
	DB       int `default:"0"`
}

type Mongo struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Kafka struct {
	Brokers                []string `default:"localhost:9092"`
	Username               string
	Password               string
	AllowAutoTopicCreation bool `default:"false"`
}

type Jwt struct {
	Secret            string
	AccessTokenExpire int64 `default:"3600"`
	RefreshExpire     int64 `default:"7200"`
	MultipleLogin     bool  `default:"true"`
}

type Email struct {
	Username string
	Password string
	Host     string
	Port     int
}

type Audit struct {
	GeoLite2 struct {
		Enabled bool
		Path    string `default:"GeoLite2-City.mmdb"`
	}
}

func (h *Http) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}
