package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DarylSerrano/dlsite-organizer/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func initializeDB(path string) string {
	var databasePath = filepath.Join(path, "data.db")
	if !fileExists(databasePath) {
		log.Print("Database doesnt exists, creating", databasePath)
		_, err := os.Create(databasePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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

func filterCircleTest(basepath string, db *sql.DB) {
	circle := CircleDB{ID: 38835, Name: "みやぢ屋"}
	filterpath := filepath.Join(basepath, "Circles")
	filterByCircle(db, circle, filterpath)
}

func filterSfwTest(basepath string, db *sql.DB) {
	folderSfw := filepath.Join(basepath, "SFW")
	filterBySfw(db, true, folderSfw)
}

func filterNsfwTest(basepath string, db *sql.DB) {
	folderNsfw := filepath.Join(basepath, "NSFW")
	filterBySfw(db, false, folderNsfw)
}

func filterTagTest(basepath string, db *sql.DB) {
	tag := TagDB{ID: 2, Name: "耳舐め"}
	folderTag := filepath.Join(basepath, "Tags")
	filterByTag(db, tag, folderTag)
}

func filterVATest(basepath string, db *sql.DB) {
	va := VoiceActorDB{ID: 1, Name: "みやぢ"}
	folderVA := filepath.Join(basepath, "VAs")
	filterByVoiceActor(db, va, folderVA)
}

func main() {
	// basepath := "./testdata"
	// databasePath := initializeDB(basepath)
	// db, err := openDB(databasePath)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer db.Close()
	// scanFiles(basepath, db)

	// filterCircleTest(basepath, db)
	// filterNsfwTest(basepath, db)
	// filterSfwTest(basepath, db)
	// filterTagTest(basepath, db)
	// filterVATest(basepath, db)

	err := cmd.Execute()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
