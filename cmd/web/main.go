package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/labstack/echo/v4"

	"formy.fprzg.net/internal/controllers"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db *sql.DB
	e  *echo.Echo
}

type Config struct {
	port  string
	dbDir string
	env   string
}

var (
	app = App{}
	cfg = Config{}
)

func main() {
	flag.StringVar(&cfg.port, "port", ":3000", "API server port.")
	flag.StringVar(&cfg.dbDir, "dbDir", "./app.db", "Database directory.")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.Parse()

	if cfg.env != "development" && cfg.env != "staging" && cfg.env != "production" {
		panic(fmt.Errorf("invalid environment: '%s'", cfg.env))
	}

	db, err := sql.Open("sqlite3", cfg.dbDir)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := controllers.GetRouter(cfg.env)

	app = App{
		db: db,
		e:  e,
	}

	app.e.Logger.Fatal(e.Start(cfg.port))
}
