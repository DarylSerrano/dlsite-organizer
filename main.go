package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal(err)
	}
	return !info.IsDir()
}

func initializeDB(path string) string {
	var dataPath = filepath.Join(path, "data")
	err := os.MkdirAll(dataPath, 0755)
	if err != nil {
		log.Panic(err)
	}
	var databasePath = filepath.Join(dataPath, "data.db")
	if !fileExists(databasePath) {
		log.Print("Database doesnt exists, creating", dataPath)
		_, err = os.Create(databasePath)
		if err != nil {
			log.Panic(err)
		}
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
	basepath := "./testdata"
	databasePath := initializeDB(basepath)
	db, err := openDB(databasePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	scanFiles(basepath, db)
	circle := CircleDB{ID: 38835, Name: "みやぢ屋"}
	filterpath := filepath.Join(basepath, "Circles")
	filterByCircle(db, circle, filterpath)
}
