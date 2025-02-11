// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/kochabonline/kcloud/apps/system/account"
	"github.com/kochabonline/kcloud/apps/system/auth/google"
	"github.com/kochabonline/kcloud/apps/system/auth/jwt"
	"github.com/kochabonline/kcloud/apps/system/menu"
	"github.com/kochabonline/kcloud/apps/system/notifier/channal"
	"github.com/kochabonline/kcloud/apps/system/notifier/message"
	"github.com/kochabonline/kcloud/apps/system/role"
	"github.com/kochabonline/kcloud/apps/system/role/bindaccount"
	"github.com/kochabonline/kcloud/apps/system/role/bindmenu"
	"github.com/kochabonline/kcloud/apps/system/security/audit"
	"github.com/kochabonline/kcloud/apps/system/security/captcha"
	"github.com/kochabonline/kcloud/apps/system/security/device"
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kcloud/internal/migrate"
	"github.com/kochabonline/kit/app"
	jwt2 "github.com/kochabonline/kit/auth/jwt"
	"github.com/kochabonline/kit/log"
	"github.com/kochabonline/kit/store/kafka"
	mongo2 "github.com/kochabonline/kit/store/mongo"
	"github.com/kochabonline/kit/store/mysql"
	redis2 "github.com/kochabonline/kit/store/redis"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func initializeMigrate(config2 *config.Config) (*migrate.Migrate, func(), error) {
	db, cleanup, err := initializeMySql(config2)
	if err != nil {
		return nil, nil, err
	}
	client, cleanup2, err := initializeMongo(config2)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	migrateMigrate := migrate.NewMigrate(db, client)
	return migrateMigrate, func() {
		cleanup2()
		cleanup()
	}, nil
}

func initializeApp(config2 *config.Config, log2 log.Helper) (*app.App, func(), error) {
	client, cleanup, err := initializeRedis(config2)
	if err != nil {
		return nil, nil, err
	}
	jwtRedis, cleanup2, err := initializeJwtRedis(config2, client)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	db, cleanup3, err := initializeMySql(config2)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	repository := account.NewRepository(db, client)
	controller := account.NewController(repository, log2)
	jwtController := jwt.NewController(jwtRedis, controller, log2)
	mongoClient, cleanup4, err := initializeMongo(config2)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	kafka, cleanup5, err := initializeKafka(config2)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	messageRepository := message.NewRepository(mongoClient, kafka, log2)
	channalRepository := channal.NewRepository(mongoClient)
	channalController := channal.NewController(channalRepository, log2)
	messageController := message.NewController(messageRepository, channalController, log2)
	queue := message.NewQueue(messageController, messageRepository, client, log2)
	handler := jwt.NewHandler(jwtController)
	googleRepository := google.NewRepository(db, client)
	googleController := google.NewController(config2, googleRepository, log2)
	googleHandler := google.NewHandler(googleController)
	accountHandler := account.NewHandler(controller)
	roleRepository := role.NewRepository(db)
	roleController := role.NewController(roleRepository, log2)
	roleHandler := role.NewHandler(roleController)
	bindaccountRepository := bindaccount.NewRepository(db)
	bindaccountController := bindaccount.NewController(bindaccountRepository, log2)
	bindaccountHandler := bindaccount.NewHandler(bindaccountController)
	bindmenuRepository := bindmenu.NewRepository(db)
	bindmenuController := bindmenu.NewController(bindmenuRepository, log2)
	bindmenuHandler := bindmenu.NewHandler(bindmenuController)
	channalHandler := channal.NewHandler(channalController)
	messageHandler := message.NewHandler(messageController)
	captchaRepository := captcha.NewRepository(client)
	captchaController := captcha.NewController(controller, captchaRepository, config2, log2)
	captchaHandler := captcha.NewHandler(captchaController)
	deviceRepository := device.NewRepository(db)
	deviceController := device.NewController(deviceRepository, log2)
	deviceHandler := device.NewHandler(deviceController)
	auditRepository := audit.NewRepository(db)
	auditController := audit.NewController(auditRepository, log2)
	auditHandler := audit.NewHandler(auditController)
	menuRepository := menu.NewRepository(db)
	menuController := menu.NewController(menuRepository, log2)
	menuHandler := menu.NewHandler(menuController)
	appApp := newApp(config2, log2, jwtController, queue, handler, googleHandler, accountHandler, roleHandler, bindaccountHandler, bindmenuHandler, channalHandler, messageHandler, captchaHandler, deviceHandler, auditHandler, menuHandler)
	return appApp, func() {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

func initializeMySql(config2 *config.Config) (*gorm.DB, func(), error) {
	m, err := mysql.New(&mysql.Config{
		Host:     config2.Mysql.Host,
		Port:     config2.Mysql.Port,
		User:     config2.Mysql.User,
		Password: config2.Mysql.Password,
		Database: config2.Mysql.Database,
		Level:    config2.Mysql.Level,
	})
	if err != nil {
		return nil, nil, err
	}

	return m.Client, func() { m.Close() }, nil
}

func initializeRedis(config2 *config.Config) (*redis.Client, func(), error) {
	r, err := redis2.NewClient(&redis2.Config{
		Host:     config2.Redis.Host,
		Port:     config2.Redis.Port,
		Password: config2.Redis.Password,
		DB:       config2.Redis.DB,
	})
	if err != nil {
		return nil, nil, err
	}

	return r.Client, func() { r.Close() }, nil
}

func initializeMongo(config2 *config.Config) (*mongo.Client, func(), error) {
	m, err := mongo2.New(&mongo2.Config{
		Host:     config2.Mongo.Host,
		Port:     config2.Mongo.Port,
		User:     config2.Mongo.User,
		Password: config2.Mongo.Password,
	})
	if err != nil {
		return nil, nil, err
	}

	return m.Client, func() { m.Close() }, nil
}

func initializeKafka(config2 *config.Config) (*kafka.Kafka, func(), error) {
	k, err := kafka.New(&kafka.Config{
		Brokers:  config2.Kafka.Brokers,
		Username: config2.Kafka.Username,
		Password: config2.Kafka.Password,
	})
	if err != nil {
		return nil, nil, err
	}

	return k, func() { k.Close() }, nil
}

var dbProviderSet = wire.NewSet(initializeMySql, initializeRedis, initializeMongo, initializeKafka)

func initializeJwtRedis(config2 *config.Config, r *redis.Client) (*jwt2.JwtRedis, func(), error) {
	j, err := jwt2.New(&jwt2.Config{
		Secret:        config2.Jwt.Secret,
		Expire:        config2.Jwt.AccessTokenExpire,
		RefreshExpire: config2.Jwt.RefreshExpire,
	})
	if err != nil {
		return nil, nil, err
	}
	jwtRedis := jwt2.NewJwtRedis(j, r, jwt2.WithMultipleLogin(config2.Jwt.MultipleLogin))

	return jwtRedis, func() {}, nil
}

var jwtProviderSet = wire.NewSet(initializeJwtRedis)
