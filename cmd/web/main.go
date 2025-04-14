package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"formy.fprzg.net/internal/types"
	"formy.fprzg.net/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var cfg types.AppConfig
	flag.StringVar(&cfg.Port, "port", ":3000", "API server port.")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (testing | development | staging | production)")

	flag.StringVar(&cfg.DBDir, "db-dir", "./app.db", "Database directory.")

	flag.StringVar(&cfg.JWTSecret, "jwt-secret", "some-secret-key", "JWT secret key.")

	// rate-limiter config
	// smtp server config

	flag.Parse()

	if cfg.Env != "development" && cfg.Env != "staging" && cfg.Env != "production" {
		panic(fmt.Errorf("invalid environment: '%s'", cfg.Env))
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(wd)

	var db *sql.DB
	//var err error
	if cfg.Env == "development" {
		db, err = utils.NewTestDB()
	} else {
		db, err = sql.Open("sqlite3", cfg.DBDir)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app, err := NewServer(cfg, db)
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Serve(); err != nil {
		log.Fatal(err)
	}
}
