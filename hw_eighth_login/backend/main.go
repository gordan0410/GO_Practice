package main

import (
	"hw_eighth_login/backend/models"
	"hw_eighth_login/backend/views"

	"github.com/rs/zerolog/log"
)



func main() {
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// errr := errors.New("Zuolar bugs")
	// log.Warn().Caller().Err(errr).Str("aa", "bb").Int("cc", 123).Msg("My Log")
	// log.Info().Str("aa", "bb").Int("cc", 123).Msg("My Log")
	// log.Error().Str("aa", "bb").Int("cc", 123).Msg("My Log")
	// log.Debug().Str("aa", "bb").Int("cc", 123).Msg("My Log")
	// log.Warn().Caller().Str("msg", "table is not exist").Str("func", "migrator.HasTable(&Account{})").Msg("DB")

	db_conn, err := models.Load_database("./config.json")
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "models.Load_database").Msg("DB")
	}
	views.Set_router(db_conn)
	models.Leave_database(db_conn)
}
