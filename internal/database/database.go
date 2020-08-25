package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/DarylSerrano/dlsite-organizer/fetcher"
	_ "github.com/mattn/go-sqlite3"
)

// A WorkDB contains information about a DLSite work on the DB
type WorkDB struct {
	id       int
	sfw      bool
	name     string
	filepath string
	circleID int
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func OpenDB(path string) (*sql.DB, error) {
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

func getTagID(db *sql.DB, tagName string) (int, error) {
	var ID int
	err := db.QueryRow("SELECT ID FROM Tags WHERE Name = ?", tagName).Scan(&ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, &errorString{s: "Not found"}
		}

		log.Panic(err)
	}
	return ID, nil
}

func getVoiceActorID(db *sql.DB, voiceActorName string) (int, error) {
	var ID int
	err := db.QueryRow("SELECT ID FROM VoiceActors WHERE Name = ?", voiceActorName).Scan(&ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, &errorString{s: "Not found"}
		}

		log.Panic(err)
	}
	return ID, nil
}

func postTag(db *sql.DB, tagName string) {
	_, err := db.Exec(`INSERT INTO Tags(Name) VALUES(?)`, tagName)
	if err != nil {
		log.Fatal(err)
	}
}

func postCircle(db *sql.DB, id string, name string) {
	log.Print("create circle ", id, name)
	_, err := db.Exec(`INSERT INTO Circles(ID, Name) VALUES(?, ?)`, id, name)
	if err != nil {
		log.Fatal(err)
	}
}

func postVoiceActor(db *sql.DB, name string) {
	_, err := db.Exec(`INSERT INTO VoiceActors(Name) VALUES(?)`, name)
	if err != nil {
		log.Fatal(err)
	}
}

func postWorkTag(db *sql.DB, workID int, tagID int) {
	_, err := db.Exec(`INSERT INTO WorkTag(WorkID, TagID) VALUES(?, ?)`, workID, tagID)
	if err != nil {
		log.Fatal(err)
	}
}

func postWorkVoiceActor(db *sql.DB, workID int, voiceActorID int) {
	_, err := db.Exec(`INSERT INTO WorkVoiceActor(WorkID, VoiceActorID) VALUES(?, ?)`,
		workID, voiceActorID)
	if err != nil {
		log.Fatal(err)
	}
}

func postWork(db *sql.DB, work fetcher.Work, filepath string) {
	log.Print("create work ", work, filepath)

	var isSfw int
	if work.SFW {
		isSfw = 1
	} else {
		isSfw = 0
	}
	_, err := db.Exec(`INSERT INTO Works(ID, sfw, Name, CircleID, Filepath) VALUES(?, ?, ?, ?, ?)`,
		work.ID, isSfw, work.Name, work.Circle.ID, filepath)
	if err != nil {
		log.Fatal(err)
	}
}

func updateWorkFilepath(db *sql.DB, work fetcher.Work, filepath string) {
	_, err := db.Exec(`UPDATE Works SET Filepath = ? WHERE ID = ?`, filepath, work.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func rowExists(db *sql.DB, tableName string, name string) bool {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT * FROM %s WHERE Name = '%s')`, tableName, name)
	log.Print("row exists of: ", query)
	var exists bool
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Panic(err)
	}
	return exists
}

func workExists(db *sql.DB, id int) bool {
	var name string
	err := db.QueryRow("SELECT Name FROM Works WHERE ID = ?", id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Panic(err)
	}
	return true
}

func circleIDExists(db *sql.DB, id int) bool {
	var name string
	err := db.QueryRow("SELECT Name FROM Circles WHERE ID = ?", id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Panic(err)
	}
	return true
}

func tagExists(db *sql.DB, name string) bool {
	return rowExists(db, "Tags", name)
}

func circleNameExists(db *sql.DB, name string) bool {
	return rowExists(db, "Circles", name)
}

func voiceActorExists(db *sql.DB, name string) bool {
	return rowExists(db, "VoiceActors", name)
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Tags
	(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Name Text NOT NULL UNIQUE,
		EnglishName Text
	);
	
	CREATE TABLE IF NOT EXISTS VoiceActors
	(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Name Text NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS Circles
	(
		ID INTEGER PRIMARY KEY,
		Name Text NOT NULL
	);
	
	-- Relationship tables
	CREATE TABLE IF NOT EXISTS Works
	(
		ID INTEGER PRIMARY KEY,
		sfw INTEGER NOT NULL,
		Name TEXT NOT NULL,
		Filepath TEXT,
		CircleID INTEGER,
	
		FOREIGN KEY (CircleID) REFERENCES Circles(ID)
	);
	
	CREATE TABLE IF NOT EXISTS WorkTag
	(
		WorkID INTEGER,
		TagID INTEGER,
	
		FOREIGN KEY (WorkID) REFERENCES Works(ID),
		FOREIGN KEY (TagID) REFERENCES Tags(ID)
	);
	
	CREATE TABLE IF NOT EXISTS WorkVoiceActor
	(
		WorkID INTEGER,
		VoiceActorID INTEGER,
	
		FOREIGN KEY (WorkID) REFERENCES Works(ID),
		FOREIGN KEY (VoiceActorID) REFERENCES VoiceActors(ID)
	);`)

	return err
}

func SaveWork(db *sql.DB, work fetcher.Work, filepath string) error {
	log.Print("Saving ", work)
	circleID, err := strconv.ParseInt(work.Circle.ID, 10, 32)
	if err != nil {
		return err
	}
	workID, err := strconv.ParseInt(work.ID, 10, 32)
	if err != nil {
		return err
	}
	if !circleIDExists(db, int(circleID)) {
		postCircle(db, work.Circle.ID, work.Circle.Name)
	}

	if workExists(db, int(workID)) {
		log.Print("WOrk exists, update")
		updateWorkFilepath(db, work, filepath)
	} else {
		postWork(db, work, filepath)
		// Save tags and create relationships
		for _, tagName := range work.Tags {
			if !tagExists(db, tagName) {
				postTag(db, tagName)
			}
			tagID, err := getTagID(db, tagName)
			if err != nil {
				return err
			}
			postWorkTag(db, int(workID), tagID)
		}

		// Save voice actors and create relationships
		for _, voiceActorName := range work.VoiceActors {
			if !voiceActorExists(db, voiceActorName) {
				postVoiceActor(db, voiceActorName)
			}
			voiceActorID, err := getVoiceActorID(db, voiceActorName)
			if err != nil {
				return err
			}
			postWorkVoiceActor(db, int(workID), voiceActorID)
		}
	}

	return nil
}
