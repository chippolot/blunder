package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "blunderbuddy",
		Usage: "generate a comical story of a misundertanding!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "themesPath",
				Aliases:  []string{"t"},
				Value:    "",
				Usage:    "Path to theme strings",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "stylesPath",
				Aliases:  []string{"s"},
				Value:    "",
				Usage:    "Path to style strings",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "modifiersPath",
				Aliases:  []string{"m"},
				Value:    "",
				Usage:    "Path to modifier strings",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dbPath",
				Aliases:  []string{"o"},
				Value:    "",
				Usage:    "Path to database file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "resetDB",
				Aliases:  []string{"r"},
				Value:    false,
				Usage:    "If true, database will be cleared before importing new strings",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			dbPath := ctx.String("dbPath")
			themesPath := ctx.String("themesPath")
			stylesPath := ctx.String("stylesPath")
			modifiersPath := ctx.String("modifiersPath")
			resetDB := ctx.Bool("resetDB")
			return prepareDb(dbPath, themesPath, stylesPath, modifiersPath, resetDB)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func prepareDb(dbPath, themesPath, stylesPath, modifiersPath string, resetDB bool) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if resetDB {
		tables := []string{"Themes", "Styles", "Modifiers"}
		fmt.Println("Resetting db...")
		for _, table := range tables {
			_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
			if err != nil {
				log.Fatalf("Failed to delete rows from %s: %v", table, err)
			}
			fmt.Printf("All rows deleted from %s table.\n", table)
		}
	}

	if err := insertStringsFromFile(db, "Themes", "Theme", themesPath); err != nil {
		return err
	}
	if err := insertStringsFromFile(db, "Styles", "Style", stylesPath); err != nil {
		return err
	}
	if err := insertStringsFromFile(db, "Modifiers", "Modifier", modifiersPath); err != nil {
		return err
	}

	fmt.Println("Finished preparing database.")
	return nil
}

func insertStringsFromFile(db *sql.DB, table string, column string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare the statement for inserting data, including a timestamp
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (%s, Timestamp) VALUES (?, ?)", table, column))
	if err != nil {
		return err
	}
	defer stmt.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		timestamp := time.Now().UTC()

		_, err := stmt.Exec(line, timestamp)
		if err != nil {
			return err
		}
		count += 1
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("Successfully inserted %v entries into table %s.\n", count, table)
	return nil
}
