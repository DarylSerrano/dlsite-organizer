package filter

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DarylSerrano/dlsite-organizer/internal/filehandler"
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
		works = append(works, work)
	}
	return works
}

func scanTags(rows *sql.Rows) []TagDB {
	var tags []TagDB
	for rows.Next() {
		var tag TagDB
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, tag)
	}
	return tags
}

func scanCircles(rows *sql.Rows) []CircleDB {
	var circles []CircleDB
	for rows.Next() {
		var circle CircleDB
		err := rows.Scan(&circle.ID, &circle.Name)
		if err != nil {
			log.Fatal(err)
		}
		circles = append(circles, circle)
	}
	return circles
}

func scanVoiceActors(rows *sql.Rows) []VoiceActorDB {
	var vas []VoiceActorDB
	for rows.Next() {
		var va VoiceActorDB
		err := rows.Scan(&va.ID, &va.Name)
		if err != nil {
			log.Fatal(err)
		}
		vas = append(vas, va)
	}
	return vas
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
	var works []WorkFilterResult
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
	works = scanWorksFiltered(rows)

	return works
}

func getFilteredWorksByTag(db *sql.DB, tagID int) []WorkFilterResult {
	var works []WorkFilterResult
	rows, err := db.Query(
		`SELECT w.ID, w.Name, w.Filepath FROM Works w 
			INNER JOIN WorkTag wt ON w.ID = wt.WorkID
			WHERE wt.TagID = ?`, tagID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	works = scanWorksFiltered(rows)

	return works
}

func getFilteredWorksByVoiceActor(db *sql.DB, voiceActorID int) []WorkFilterResult {
	var works []WorkFilterResult
	rows, err := db.Query(
		`SELECT w.ID, w.Name, w.Filepath FROM Works w 
			INNER JOIN WorkVoiceActor wva ON w.ID = wva.WorkID
			WHERE wva.VoiceActorID = ?`, voiceActorID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	works = scanWorksFiltered(rows)

	return works
}

// GetAllTags returns all tags stored on DB
func GetAllTags(db *sql.DB) []TagDB {
	rows, err := db.Query(`SELECT ID, NAME FROM Tags`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var tags = scanTags(rows)
	return tags
}

// GetAllCircles returns all circles stored on DB
func GetAllCircles(db *sql.DB) []CircleDB {
	rows, err := db.Query(`SELECT ID, NAME FROM Circles`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var circles = scanCircles(rows)
	return circles
}

// GetAllVoiceActors returns all Voice Actors stored opn DB
func GetAllVoiceActors(db *sql.DB) []VoiceActorDB {
	rows, err := db.Query(`SELECT ID, Name FROM VoiceActors`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var vas = scanVoiceActors(rows)
	return vas
}

// ByCircle filters works by Circle and create symbolic links for the files of the works filtered
func ByCircle(db *sql.DB, circle CircleDB, basepath string) {
	works := getFilteredWorskByCircle(db, circle.ID)
	circleFolder := filepath.Join(basepath, circle.Name)
	os.MkdirAll(circleFolder, 0755)

	for _, work := range works {
		filename := fmt.Sprint("RJ", work.ID)
		newName := filepath.Join(circleFolder, filename)
		filehandler.CreateSymlink(work.filepath, newName)
	}
}

// BySfw filters works by SFW or NSFW and create symbolic links for the files of the works filtered
func BySfw(db *sql.DB, isSfw bool, basepath string) {
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
		filehandler.CreateSymlink(work.filepath, newName)
	}
}

// ByTag filters works by Tag and create symbolic links for the files of the works filtered
func ByTag(db *sql.DB, tag TagDB, basepath string) {
	works := getFilteredWorksByTag(db, tag.ID)
	tagFolder := filepath.Join(basepath, tag.Name)
	os.MkdirAll(tagFolder, 0755)
	for _, work := range works {
		filename := fmt.Sprint("RJ", work.ID)
		newName := filepath.Join(tagFolder, filename)
		filehandler.CreateSymlink(work.filepath, newName)
	}
}

// ByVoiceActor filters works by Voice Actor and create symbolic links for the files of the works filtered
func ByVoiceActor(db *sql.DB, va VoiceActorDB, basepath string) {
	works := getFilteredWorksByVoiceActor(db, va.ID)
	vaFolder := filepath.Join(basepath, va.Name)
	os.MkdirAll(vaFolder, 0755)
	for _, work := range works {
		filename := fmt.Sprint("RJ", work.ID)
		newName := filepath.Join(vaFolder, filename)
		filehandler.CreateSymlink(work.filepath, newName)
	}
}
