package main

import (
	"database/sql"
	"os"
	"weight-tracker/pkg/api"
	"weight-tracker/pkg/app"
	"weight-tracker/pkg/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

// func run will be responsible for setting up db connections, routers etc
func run() error {
	// I'm used to working with postgres, but feel free to use any db you like. You just have to change the driver
	// I'm not going to cover how to create a database here but create a database
	// and call it something along the lines of "weight tracker"
	// connectionString := "postgres://postgres:postgres@localhost/**NAME-OF-YOUR-DATABASE-HERE**?sslmode=disable"
	connectionString := "root:root@tcp(localhost:3306)/structure_test"

	// setup database connection
	db, err := setupDatabase(connectionString)

	if err != nil {
		log.Error().Caller().Err(err).Msg("err")
		return err
	}

	// create storage dependency
	storage := repository.NewStorage(db)

	// run migrations
	// note that we are passing the connectionString again here. This is so
	// we can easily run migrations against another database, say a test version,
	// for our integration- and end-to-end tests
	err = storage.RunMigrations(connectionString)

	if err != nil {
		log.Error().Caller().Err(err).Msg("err")
		return err
	}

	// create router dependecy
	router := gin.Default()
	router.Use(cors.Default())

	// create user service
	userService := api.NewUserService(storage)

	// create weight service
	weightService := api.NewWeightService(storage)
	// fmt.Println(userService, weightService)
	server := app.NewServer(router, userService, weightService)

	// start the server
	err = server.Run()

	if err != nil {
		log.Error().Caller().Err(err).Msg("err")
		return err
	}

	return nil
}

func setupDatabase(connString string) (*sql.DB, error) {
	// change "postgres" for whatever supported database you want to use
	db, err := sql.Open("mysql", connString)

	if err != nil {
		return nil, err
	}

	// ping the DB to ensure that it is connected
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
