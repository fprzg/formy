package main

import (
	"database/sql"
	"flag"
	"fmt"

	"formy.fprzg.net/internal/controllers"
	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"

	_ "github.com/mattn/go-sqlite3"
)

var (
	cfg = types.AppConfig{}
)

func main() {
	flag.StringVar(&cfg.Port, "port", ":3000", "API server port.")
	flag.StringVar(&cfg.DBDir, "dbDir", "./app.db", "Database directory.")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development | staging | production)")
	flag.Parse()

	if cfg.Env != "development" && cfg.Env != "staging" && cfg.Env != "production" {
		panic(fmt.Errorf("invalid environment: '%s'", cfg.Env))
	}

	var m *models.Models
	var err error

	if cfg.Env == "development" {
		m, err = models.GetTestModels()
		if err != nil {
			panic(err)
		}
	} else {
		db, err := sql.Open("sqlite3", cfg.DBDir)
		if err != nil {
			panic(err)
		}

		m = models.GetModels(db)
	}

	c := controllers.GetControllers(m)

	c.Start(cfg)
}
