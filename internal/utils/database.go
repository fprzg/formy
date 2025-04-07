package utils

import (
	"database/sql"
	"fmt"
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

func ExecuteSqlStmt(db *sql.DB, stmt string, args ...any) (int64, error) {
	result, err := db.Exec(stmt, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func SetupTestDB() (*sql.DB, error) {
	return SetupDB(":memory:")
}

func SetupDB(dbPath string) (*sql.DB, error) {
	const migrationsDir = "../../migrations"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	migrations, err := LoadMigrations(migrationsDir)
	if err != nil {
		return nil, err
	}

	err = ApplyMigrations(db, migrations, GetLatestVersion(migrations))
	if err != nil {
		return nil, err
	}

	return db, nil
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

func ApplyMigrations(db *sql.DB, migrations []Migration, target int) error {
	// Going up or down?
	direction := "up"
	current := 0
	start, end := current+1, target
	if target < current {
		direction = "down"
		start, end = current, target+1
	}

	for _, m := range migrations {
		if (direction == "up" && m.Type == "up" && m.Version >= start && m.Version <= end) ||
			(direction == "down" && m.Type == "down" && m.Version < start && m.Version >= end) {
			content, err := os.ReadFile(m.Path)
			if err != nil {
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return fmt.Errorf("failed to apply migration %d (%s): %v", m.Version, m.Type, err)
			}
		}
	}

	return nil
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

	return nil
}

/*
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
*/
