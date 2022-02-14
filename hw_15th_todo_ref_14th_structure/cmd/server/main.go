package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"todolist/internal/api"
	"todolist/internal/app"
	"todolist/internal/repository"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type Configs struct {
	Jwt struct {
		JwtKey    string `json:"Jwt_key"`
		JwtMaxage int    `json:"Jwt_maxage"`
	} `json:"jwt"`
	Mysql struct {
		Username     string `json:"Username"`
		Password     string `json:"Password"`
		Protocol     string `json:"Protocol"`
		Host         string `json:"Host"`
		Port         string `json:"Port"`
		Database     string `json:"Database"`
		MaxLifetime  int    `json:"Max_lifetime"`
		MaxOpenconns int    `json:"Max_openconns"`
		MaxIdleconns int    `json:"Max_idleconns"`
	} `json:"mysql"`
	Redis struct {
		Size     int    `json:"Size"`
		Network  string `json:"Network"`
		Address  string `json:"Address"`
		Password string `json:"Password"`
	} `json:"redis"`
	Session struct {
		SessionName   string `json:"Session_name"`
		SessionPrefix string `json:"Session_prefix"`
		SessionKey    string `json:"Session_key"`
		SessionMaxage int    `json:"Session_maxage"`
	} `json:"session"`
}

func main() {
	var configAddr string
	var migrateCtl bool
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			if os.Args[2] != "" {
				configAddr = os.Args[2]
			} else {
				configAddr = "./configs/config.json"
			}
		case "migrate":
			migrateCtl = true
			configAddr = "./configs/config.json"
		default:
			configAddr = "./configs/config.json"
		}
	} else {
		configAddr = "./configs/config.json"
	}
	if err := run(configAddr, migrateCtl); err != nil {
		os.Exit(1)
	}
}

func run(configAddr string, migrateCtl bool) error {
	// load configs
	conf, err := loadConfig(configAddr)
	if err != nil {
		log.Error().Caller().Err(err).Msg("config")
		return err
	}

	// mysql connection string
	connectionString := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.Mysql.Username, conf.Mysql.Password, conf.Mysql.Protocol, conf.Mysql.Host, conf.Mysql.Port, conf.Mysql.Database)
	// set db
	db, err := setupDatabase(connectionString, conf)
	if err != nil {
		log.Error().Caller().Err(err).Msg("DB")
		return err
	}
	defer db.Close()
	// new db service
	storage := repository.NewStorage(db)
	// migrate
	if migrateCtl {
		err := storage.RunMigration()
		if err != nil {
			log.Error().Caller().Err(err).Msg("DB")
			return err
		}
	}

	// set ctx
	ctx := context.Background()
	// set redis client
	redis, err := setupRedis(conf, ctx)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Redis")
		return err
	}
	// new redis servic
	redisStorage := repository.NewRedisStorage(redis, ctx)

	// gin
	router := gin.Default()

	// new todolist service
	todolistService := api.NewTodolistService(storage)

	// new session
	sessionService := api.NewSessionService(redisStorage, conf.Session.SessionName, conf.Session.SessionKey, conf.Session.SessionPrefix, conf.Session.SessionMaxage)

	// new login service
	loginService := api.NewLoginService(storage, sessionService)

	// new Server
	server := app.NewServer(router, todolistService, loginService)

	err = server.Run()
	if err != nil {
		log.Warn().Caller().Err(err).Msg("server")
		return err
	}

	return nil

}

func loadConfig(addr string) (*Configs, error) {
	file, err := os.Open(addr)
	if err != nil {
		return nil, err
	}
	var c *Configs
	json.NewDecoder(file).Decode(&c)
	return c, nil
}

func setupDatabase(connString string, c *Configs) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	db.DB().SetConnMaxLifetime(time.Duration(c.Mysql.MaxLifetime))
	db.DB().SetMaxOpenConns(c.Mysql.MaxOpenconns)
	db.DB().SetMaxIdleConns(c.Mysql.MaxIdleconns)
	return db, nil
}

func setupRedis(c *Configs, ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Address,
		Password: "",
		DB:       0,
		Network:  c.Redis.Network,
		PoolSize: c.Redis.Size,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
