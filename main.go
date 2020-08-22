package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func initializeDB(path string) string {
	var dataPath = filepath.Join(path, "data")
	err := os.MkdirAll(dataPath, 0755)
	if err != nil {
		log.Panic(err)
	}
	var databasePath = filepath.Join(dataPath, "data.db")
	_, err = os.Create(databasePath)
	if err != nil {
		log.Panic(err)
	}

	return databasePath
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	err = createTables(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	databasePath := initializeDB(".")
	db, err := openDB(databasePath)
	if err != nil {
		log.Panic(err)
	}
	scanFiles(".", db)
}
