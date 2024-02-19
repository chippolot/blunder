package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chippolot/blunders/internal/blunder"
)

type SQLiteDataProvider struct {
	db *sql.DB
}

func MakeSQLiteDataProvider() *SQLiteDataProvider {
	var err error
	db, err := sql.Open("sqlite3", "blunder.db")
	if err != nil {
		panic(err)
	}

	err = createTables(db)
	if err != nil {
		panic(err)
	}

	err = createIndexes(db)
	if err != nil {
		panic(err)
	}

	dataProvider := &SQLiteDataProvider{
		db: db,
	}

	return dataProvider
}

func (f *SQLiteDataProvider) GetRandomString(dataType blunder.StoryDataType) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func createIndexes(db *sql.DB) error {
	var err error

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_stories_timestamp ON Stories (Timestamp DESC);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return err
	}

	return nil
}

func createTables(db *sql.DB) error {
	var err error

	createTableSQL := `CREATE TABLE IF NOT EXISTS Stories (
		"Id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"Story" TEXT,
		"Prompt" TEXT,
		"Timestamp" DATETIME
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	createTableSQL = `CREATE TABLE IF NOT EXISTS Themes (
		"Id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"Theme" TEXT,
		"Timestamp" DATETIME
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	createTableSQL = `CREATE TABLE IF NOT EXISTS Styles (
		"Id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"Style" TEXT,
		"Timestamp" DATETIME
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	createTableSQL = `CREATE TABLE IF NOT EXISTS Modifiers (
		"Id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"Modifier" TEXT,
		"Timestamp" DATETIME
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}
