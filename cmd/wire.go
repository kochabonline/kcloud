//go:build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/kochabonline/kcloud/apps/system"
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kcloud/internal/migrate"
	"github.com/kochabonline/kit/app"
	"github.com/kochabonline/kit/auth/jwt"
	"github.com/kochabonline/kit/log"
	kkafka "github.com/kochabonline/kit/store/kafka"
	kmongo "github.com/kochabonline/kit/store/mongo"
	kmysql "github.com/kochabonline/kit/store/mysql"
	kredis "github.com/kochabonline/kit/store/redis"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func initializeMySql(config *config.Config) (*gorm.DB, func(), error) {
	m, err := kmysql.New(&kmysql.Config{
		Host:     config.Mysql.Host,
		Port:     config.Mysql.Port,
		User:     config.Mysql.User,
		Password: config.Mysql.Password,
		Database: config.Mysql.Database,
		Level:    config.Mysql.Level,
	})
	if err != nil {
		return nil, nil, err
	}

	return m.Client, func() { m.Close() }, nil
}

func initializeRedis(config *config.Config) (*redis.Client, func(), error) {
	r, err := kredis.NewClient(&kredis.Config{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err != nil {
		return nil, nil, err
	}

	return r.Client, func() { r.Close() }, nil
}

func initializeMongo(config *config.Config) (*mongo.Client, func(), error) {
	m, err := kmongo.New(&kmongo.Config{
		Host:     config.Mongo.Host,
		Port:     config.Mongo.Port,
		User:     config.Mongo.User,
		Password: config.Mongo.Password,
	})
	if err != nil {
		return nil, nil, err
	}

	return m.Client, func() { m.Close() }, nil
}

func initializeKafka(config *config.Config) (*kkafka.Kafka, func(), error) {
	k, err := kkafka.New(&kkafka.Config{
		Brokers:  config.Kafka.Brokers,
		Username: config.Kafka.Username,
		Password: config.Kafka.Password,
	})
	if err != nil {
		return nil, nil, err
	}

	return k, func() { k.Close() }, nil
}

var dbProviderSet = wire.NewSet(initializeMySql, initializeRedis, initializeMongo, initializeKafka)

func initializeJwtRedis(config *config.Config, r *redis.Client) (*jwt.JwtRedis, func(), error) {
	j, err := jwt.New(&jwt.Config{
		Secret:        config.Jwt.Secret,
		Expire:        config.Jwt.AccessTokenExpire,
		RefreshExpire: config.Jwt.RefreshExpire,
	})
	if err != nil {
		return nil, nil, err
	}
	jwtRedis := jwt.NewJwtRedis(j, r, jwt.WithMultipleLogin(config.Jwt.MultipleLogin))

	return jwtRedis, func() {}, nil
}

var jwtProviderSet = wire.NewSet(initializeJwtRedis)

func initializeMigrate(config *config.Config) (*migrate.Migrate, func(), error) {
	panic(wire.Build(dbProviderSet, migrate.ProviderSet))
}

func initializeApp(config *config.Config, log log.Helper) (*app.App, func(), error) {
	panic(wire.Build(dbProviderSet, jwtProviderSet, system.ProviderSet, newApp))
}
