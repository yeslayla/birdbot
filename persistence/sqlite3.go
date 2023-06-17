package persistence

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql/*
var migrationScripts embed.FS

type Sqlite3Database struct {
	db *sql.DB
}

// NewSqlite3Database creates a new SqliteDB object
func NewSqlite3Database(path string) *Sqlite3Database {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Printf("failed to create directory for db: %s", err)
			return nil
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Printf("failed to open db: %s", err)
		return nil
	}

	return &Sqlite3Database{
		db: db,
	}
}

func getMigrations() migrate.MigrationSource {
	return &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationScripts,
		Root:       "sql",
	}
}

// MigrateUp migrates the DB
func (db *Sqlite3Database) MigrateUp() error {

	n, err := migrate.Exec(db.db, "sqlite3", getMigrations(), migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to migrate: %s", err)
	}

	if n != 0 {
		log.Printf("Applied %d DB migrations!\n", n)
	}
	return nil
}

// MigrateUp destroys the DB
func (db *Sqlite3Database) MigrateDown() error {

	n, err := migrate.Exec(db.db, "sqlite3", getMigrations(), migrate.Down)
	if err != nil {
		return fmt.Errorf("failed to migrate: %s", err)
	}

	if n != 0 {
		log.Printf("Applied %d DB migrations!\n", n)
	}
	return nil
}

// GetDiscordMessage finds a discord message ID from a given local ID
func (db *Sqlite3Database) GetDiscordMessage(id string) (string, error) {

	var messageID string
	row := db.db.QueryRow("SELECT message_id FROM discord_messages WHERE id = $1", id)

	if err := row.Scan(&messageID); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get discord message from sqlite3: %s", err)
	}

	return messageID, nil
}

// SetDiscordMessage sets a discord message ID from a given local ID
func (db *Sqlite3Database) SetDiscordMessage(id string, messageID string) error {

	statement, err := db.db.Prepare("INSERT OR IGNORE INTO discord_messages (id, message_id) VALUES (?, ?)")
	if err != nil {
		return err
	}

	result, err := statement.Exec(id, messageID)
	if err != nil {
		return err
	}

	n, _ := result.RowsAffected()

	if n == 0 {
		statement, err := db.db.Prepare("UPDATE discord_messages SET message_id = (?) WHERE id = (?)")
		if err != nil {
			return err
		}

		if _, err := statement.Exec(messageID, id); err != nil {
			return err
		}

	}

	return nil
}

// GetDiscordWebhook finds a discord webhook based on a given local id
func (db *Sqlite3Database) GetDiscordWebhook(id string) (*DBDiscordWebhook, error) {

	var data DBDiscordWebhook = DBDiscordWebhook{}
	row := db.db.QueryRow("SELECT webhook_id, webhook_token FROM discord_webhooks WHERE id = $1", id)

	if err := row.Scan(&data.ID, &data.Token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get discord webhook from sqlite3: %s", err)
	}

	return &data, nil
}

// SetDiscordWebhook stores a discord webhook based on a given local id
func (db *Sqlite3Database) SetDiscordWebhook(id string, data *DBDiscordWebhook) error {

	statement, err := db.db.Prepare("INSERT OR IGNORE INTO discord_webhooks (id, webhook_id, webhook_token) VALUES (?, ?)")
	if err != nil {
		return err
	}

	result, err := statement.Exec(id, data.ID, data.Token)
	if err != nil {
		return err
	}

	n, _ := result.RowsAffected()

	if n == 0 {
		statement, err := db.db.Prepare("UPDATE discord_webhooks SET webhook_id = (?), webhook_token = (?) WHERE id = (?)")
		if err != nil {
			return err
		}

		if _, err := statement.Exec(data.ID, data.Token, id); err != nil {
			return err
		}

	}

	return nil
}
