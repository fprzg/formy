package main

import (
	"flag"
	"fmt"
	"log"

	"formy.fprzg.net/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

type cfg struct {
	MigrationsDir string
	TargetVersion int
	DBPath        string
	GenerateName  string
}

func main() {
	cfg := cfg{}
	flag.StringVar(&cfg.MigrationsDir, "dir", "", "directory containing migration scripts")
	flag.IntVar(&cfg.TargetVersion, "target", -1, "target migration version (default: latest)")
	flag.StringVar(&cfg.DBPath, "db", "./app.db", "path to the application SQLite database")
	flag.StringVar(&cfg.GenerateName, "generate", "", "generate new up/down migration files with this name")
	flag.Parse()

	if cfg.MigrationsDir == "" {
		log.Fatalf("Migration directory not defined.\n")
	}

	// Generate migration files option
	if cfg.GenerateName != "" {
		err := utils.GenerateMigrationFiles(cfg.MigrationsDir, cfg.GenerateName)
		if err != nil {
			log.Fatalf("Failed to generate migration files: %v", err)
		}
		fmt.Printf("Generated migration files for '%s'\n", cfg.GenerateName)
		return
	}

	// Apply migrations files option
	ctx, err := utils.NewMigrationCtx(cfg.DBPath, "./migration_state.db", cfg.MigrationsDir)
	if err != nil {
		log.Fatal("" + err.Error())
	}

	err = ctx.Migrate(cfg.TargetVersion)
	if err != nil {
		log.Fatalf("" + err.Error())
	}

	ctx.Close()
}
