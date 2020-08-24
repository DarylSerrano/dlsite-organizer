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

// A TagDB contains data of the tag on the db
type TagDB struct {
	ID   int
	Name string
}

// A VoiceActorDB contains data of the VA on the DB
type VoiceActorDB struct {
	ID   int
	Name string
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

func getFilteredWorskByCircle(db *sql.DB, circleID int) []WorkFilterResult {
	rows, err := db.Query("SELECT ID, Name, Filepath FROM Works WHERE CircleID = ?", circleID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var works = scanWorksFiltered(rows)
	return works
}

func getFilteredWorksBySfw(db *sql.DB, isSfw bool) []WorkFilterResult {
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

func getFilteredWorksByTag(db *sql.DB, tagID int) []WorkFilterResult {
	rows, err := db.Query(
		`SELECT w.ID, w.Name, w.Filepath FROM Works w 
			INNER JOIN WorkTag wt ON w.ID = wt.WorkID
			WHERE wt.TagID = ?`, tagID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var works = scanWorksFiltered(rows)
	return works
}

func getFilteredWorksByVoiceActor(db *sql.DB, voiceActorID int) []WorkFilterResult {
	rows, err := db.Query(
		`SELECT w.ID, w.Name, w.Filepath FROM Works w 
			INNER JOIN WorkVoiceActor wva ON w.ID = wva.WorkID
			WHERE wva.VoiceActorID = ?`, voiceActorID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var works = scanWorksFiltered(rows)
	return works
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
	works := getFilteredWorksBySfw(db, isSfw)
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

func filterByTag(db *sql.DB, tag TagDB, basepath string) {
	works := getFilteredWorksByTag(db, tag.ID)
	tagFolder := filepath.Join(basepath, tag.Name)
	os.MkdirAll(tagFolder, 0755)
	for _, work := range works {
		filename := fmt.Sprint("RJ", work.ID)
		newName := filepath.Join(tagFolder, filename)
		createSymlink(work.filepath, newName)
	}
}

func filterByVoiceActor(db *sql.DB, va VoiceActorDB, basepath string) {
	works := getFilteredWorksByVoiceActor(db, va.ID)
	vaFolder := filepath.Join(basepath, va.Name)
	os.MkdirAll(vaFolder, 0755)
	for _, work := range works {
		filename := fmt.Sprint("RJ", work.ID)
		newName := filepath.Join(vaFolder, filename)
		createSymlink(work.filepath, newName)
	}
}
