package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Migration struct {
	Version int
	Name    string
	Type    string // "up" or "down"
	Path    string
}

type MigrationCtx struct {
	AppDB          *sql.DB
	StateDB        *sql.DB
	CurrentVersion int
	MigrationsDir  string
}

// TODO(Farid): Refactor this function so we can both migrate "up" and "down".
func NewMigrationCtx(appDBPath, stateDBPath, migrationsDir string) (MigrationCtx, error) {
	stateDB, err := sql.Open("sqlite3", stateDBPath)
	if err != nil {
		return MigrationCtx{}, fmt.Errorf("Failed to open state database: " + err.Error())
	}
	currentVersion := InitializeStateDB(stateDB)

	appDB, err := sql.Open("sqlite3", appDBPath)
	if err != nil {
		return MigrationCtx{}, fmt.Errorf("Failed to open application database: " + err.Error())
	}

	ctx := MigrationCtx{
		AppDB:          appDB,
		StateDB:        stateDB,
		CurrentVersion: currentVersion,
		MigrationsDir:  migrationsDir,
	}

	return ctx, nil
}

func (ctx *MigrationCtx) Close() {
	ctx.AppDB.Close()
	ctx.StateDB.Close()
}

func (ctx *MigrationCtx) Migrate(target int) error {
	migrations, err := LoadMigrations(ctx.MigrationsDir)
	if err != nil {
		return fmt.Errorf("Failed to load migrations: " + err.Error())
	}

	if target == -1 {
		target = GetLatestVersion(migrations)
	}

	err = ApplyMigrations(ctx.AppDB, ctx.StateDB, migrations, ctx.CurrentVersion, target)
	if err != nil {
		return fmt.Errorf("Failed to apply migrations: " + err.Error())
	}

	//log.Printf("Migration completed. Current version: %d\n", target)
	return nil
}

func LoadMigrations(dir string) ([]Migration, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		rest := parts[1]
		typeParts := strings.Split(rest, ".")
		if len(typeParts) < 3 {
			continue
		}

		migrationType := typeParts[len(typeParts)-2]
		if migrationType != "up" && migrationType != "down" {
			continue
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    strings.Join(typeParts[:len(typeParts)-2], "."),
			Type:    migrationType,
			Path:    filepath.Join(dir, name),
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getLatestVersion returns the highest version number from migrations
func GetLatestVersion(migrations []Migration) int {
	if len(migrations) == 0 {
		return 0
	}
	return migrations[len(migrations)-1].Version
}

// returns the current migration version
func InitializeStateDB(stateDB *sql.DB) int {
	var tableExists bool
	err := stateDB.QueryRow(`
    SELECT EXISTS (
        SELECT 1 
        FROM sqlite_master 
        WHERE type='table' 
        AND name='migration_state'
    )
	`).Scan(&tableExists)
	if err != nil {
		log.Fatalf("Failed to check if migration table exists: %v", err)
	}

	if !tableExists {
		_, err = stateDB.Exec(`
        CREATE TABLE migration_state (
            current_version INTEGER PRIMARY KEY
        );
        INSERT INTO migration_state (current_version) VALUES (0);
    `)
		if err != nil {
			log.Fatalf("Failed to create migration table: %v", err)
		}
		//log.Println("Created migration_state table with initial version 0")
	}

	// Now get the current version
	var currentVersion int
	err = stateDB.QueryRow("SELECT current_version FROM migration_state").Scan(&currentVersion)
	if err != nil {
		log.Fatalf("Failed to get current version: %v", err)
	}
	//log.Printf("Current database version: %d", currentVersion)

	return currentVersion
}

func ApplyMigrations(appDB, stateDB *sql.DB, migrations []Migration, current, target int) error {
	if current == target {
		//log.Println("No migrations needed. Already at target version.")
		return nil
	}

	tx, err := stateDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Going up or down?
	direction := "up"
	start, end := current+1, target
	if target < current {
		direction = "down"
		start, end = current, target+1
	}

	for _, m := range migrations {
		if (direction == "up" && m.Type == "up" && m.Version >= start && m.Version <= end) ||
			(direction == "down" && m.Type == "down" && m.Version < start && m.Version >= end) {
			err := ExecuteMigration(appDB, m)
			if err != nil {
				return fmt.Errorf("failed to apply migration %d (%s): %v", m.Version, m.Type, err)
			}
			//log.Printf("Applied %s migration: %d\n", m.Type, m.Version)
		}
	}

	// Update state
	_, err = tx.Exec("UPDATE migration_state SET current_version = ?", target)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func ExecuteMigration(db *sql.DB, m Migration) error {
	content, err := os.ReadFile(m.Path)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(content))
	return err
}

func GenerateMigrationFiles(dir, name string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Load existing migrations to find the next version
	migrations, err := LoadMigrations(dir)
	if err != nil {
		return err
	}
	nextVersion := GetLatestVersion(migrations) + 1

	versionStr := fmt.Sprintf("%06d", nextVersion)

	upFile := filepath.Join(dir, fmt.Sprintf("%s_%s.up.sql", versionStr, name))
	if err := os.WriteFile(upFile, []byte("-- Up migration\n"), 0644); err != nil {
		return err
	}

	downFile := filepath.Join(dir, fmt.Sprintf("%s_%s.down.sql", versionStr, name))
	if err := os.WriteFile(downFile, []byte("-- Down migration\n"), 0644); err != nil {
		os.Remove(upFile)
		return err
	}

	//fmt.Printf("Created: %s\n", upFile)
	//fmt.Printf("Created: %s\n", downFile)
	return nil
}
