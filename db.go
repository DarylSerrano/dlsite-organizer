package main

import (
	"database/sql"
	"log"

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

func queryTest(db *sql.DB) {
	var (
		id       int
		sfw      bool
		name     string
		circleID int
	)
	rows, err := db.Query("select * from Works where ID = ?", 293003)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &sfw, &name, &circleID)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, sfw, name, circleID)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func getAllWorks(db *sql.DB) []WorkDB {
	var works []WorkDB

	rows, err := db.Query("select * from Works")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var work WorkDB
		err := rows.Scan(&work.id, &work.sfw, &work.name, &work.circleID, &work.filepath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(work.id, work.sfw, work.name, work.circleID, work.filepath)
		works = append(works, work)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return works
}

func insertTag(db *sql.DB, tagName string) {
	stmt, err := db.Prepare("INSERT INTO tags(Name) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(tagName)
	if err != nil {
		log.Fatal(err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
}

func access() {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// queryTest(db)

	// getAllWorks(db)

	// insertTag(db, "ささやき")
}

func containsWork(db *sql.DB, id int) bool {
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
