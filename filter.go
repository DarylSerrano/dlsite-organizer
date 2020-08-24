package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// A WorkFilterResult contains information of queried work filtered
type WorkFilterResult struct {
	ID       string
	Name     string
	filepath string
}

// A CircleDB represents data of the circle on the db
type CircleDB struct {
	ID   int
	Name string
}

func filterByCircle(db *sql.DB, circle CircleDB, basepath string) {
	works := getFilteredWorskByCircle(db, circle.ID)
	// Create Circle folder
	circleFolder := filepath.Join(basepath, circle.Name)
	os.MkdirAll(circleFolder, 0755)
	// Filter
	// Each work create symlink
	for _, work := range works {
		filename := fmt.Sprint("RJ", work.ID)
		newName := filepath.Join(circleFolder, filename)
		createSymlink(work.filepath, newName)
	}
}

func filterBySfw(db *sql.DB, isSfw bool, basepath string) {
	works := getFilterWorksBySfw(db, isSfw)
	var folderName string
	if isSfw {
		folderName = "sfw"
	} else {
		folderName = "nsfw"
	}
	sfwFolder := filepath.Join(basepath, folderName)
	os.MkdirAll(sfwFolder, 0755)
	for _, work := range works {
		filename := fmt.Sprint("RJ", work.ID)
		newName := filepath.Join(sfwFolder, filename)
		createSymlink(work.filepath, newName)
	}
}

func getFilteredWorskByCircle(db *sql.DB, circleID int) []WorkFilterResult {
	rows, err := db.Query("SELECT ID, Name, Filepath FROM Works WHERE CircleID = ?", circleID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var works = scanWorksFiltered(rows)
	return works
}

func scanWorksFiltered(rows *sql.Rows) []WorkFilterResult {
	var works []WorkFilterResult
	for rows.Next() {
		var work WorkFilterResult
		err := rows.Scan(&work.ID, &work.Name, &work.filepath)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("work", work)
		works = append(works, work)
	}
	return works
}

func getFilterWorksBySfw(db *sql.DB, isSfw bool) []WorkFilterResult {
	var sfw int
	if isSfw {
		sfw = 1
	} else {
		sfw = 0
	}

	rows, err := db.Query("SELECT ID, Name, Filepath FROM Works WHERE sfw = ?", sfw)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var works = scanWorksFiltered(rows)
	return works
}

func filterWorksByVoiceActor(db *sql.DB, voiceActorID int) {

}

func filterWorksByTag(db *sql.DB, tagName string) {

}
