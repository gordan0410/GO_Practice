package main

import (
	"hw_ninth/apps/login"
	"hw_ninth/apps/server"
	"hw_ninth/apps/todolist"
	"hw_ninth/models"
	"hw_ninth/tools"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	configs, err := tools.Load_configs("./configs/config.json")
	if err != nil {
		log.Error().Caller().Str("func", "tools.Load_configs(\"./configs/config.json\"").Err(err).Msg("Web")
		return
	}
	if len((os.Args)) == 2 {
		arg := os.Args[1]
		switch arg {
		case "all":
			go login.Set_router(configs)
			go todolist.Set_router(configs)
			server.Set_router()
		case "server":
			server.Set_router()
		case "login":
			login.Set_router(configs)
		case "todolist":
			todolist.Set_router(configs)
		case "migrate":
			db, err := models.Db_init(configs, false)
			if err != nil {
				log.Error().Caller().Str("func", "models.Db_init(configs, false)").Err(err).Msg("Web")
				return
			}
			err = models.Db_migrate(db)
			if err != nil {
				log.Error().Caller().Str("func", "models.Db_migrate(db)").Err(err).Msg("Web")
				return
			}
			err = models.Db_close(db)
			if err != nil {
				log.Error().Caller().Str("func", "models.Db_close(db)").Err(err).Msg("Web")
				return
			}
			log.Print("migrate success")
		default:
			log.Print("wrong argument provided, only: all, server, login, todolist, migrate.")
		}
	} else {
		// 開啟gin
		go login.Set_router(configs)
		// 開啟gin
		go todolist.Set_router(configs)
		// 開啟gin
		server.Set_router()
	}
}
