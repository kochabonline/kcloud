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
	Http  Http  `json:"http"`
	Mysql Mysql `json:"mysql"`
	Redis Redis `json:"redis"`
	Mongo Mongo `json:"mongo"`
	Kafka Kafka `json:"kafka"`
	Jwt   Jwt   `json:"jwt"`
	Email Email `json:"email"`
	Audit Audit `json:"audit"`
}

type Http struct {
	Host    string `json:"host" default:"localhost"`
	Port    int    `json:"port" default:"8080"`
	Level   string `json:"level" default:"release"`
	Swagger struct {
		Enabled bool `json:"enabled"`
	}
	Health struct {
		Enabled bool `json:"enabled"`
	}
	Metrics struct {
		Enabled bool `json:"enabled"`
	}
}

type Mysql struct {
	Host     string `json:"host" default:"localhost"`
	Port     int    `json:"port" default:"3306"`
	User     string `json:"user" default:"root"`
	Password string `json:"password"`
	Database string `json:"database"`
	Level    string `json:"level" default:"silent"`
}

type Redis struct {
	Host     string `json:"host" default:"localhost"`
	Port     int    `json:"port" default:"6379"`
	Password string `json:"password"`
	DB       int    `json:"db" default:"0"`
}

type Mongo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Kafka struct {
	Brokers                []string `json:"brokers" default:"localhost:9092"`
	Username               string   `json:"username"`
	Password               string   `json:"password"`
	AllowAutoTopicCreation bool     `json:"allowAutoTopicCreation" default:"false"`
}

type Jwt struct {
	Secret            string `json:"secret"`
	AccessTokenExpire int64  `json:"accessTokenExpire" default:"3600"`
	RefreshExpire     int64  `json:"refreshExpire" default:"7200"`
	MultipleLogin     bool   `json:"multipleLogin" default:"true"`
}

type Email struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Audit struct {
	GeoLite2 struct {
		Enabled bool   `json:"enabled"`
		Path    string `json:"path" default:"GeoLite2-City.mmdb"`
	} `json:"geoLite2"`
}

func (h *Http) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}
