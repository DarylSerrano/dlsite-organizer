package filehandler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/DarylSerrano/dlsite-organizer/fetcher"
	"github.com/DarylSerrano/dlsite-organizer/internal/database"
)

var reCode = regexp.MustCompile(`R(J|G)\d+`)

// CreateDBFile crates DB file if doesnt exists
func CreateDBFile(path string) string {
	var databasePath = filepath.Join(path, "data.db")
	if !FileExists(databasePath) {
		fmt.Printf("Database doesnt exists, creating db on: %v\n", databasePath)
		_, err := os.Create(databasePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	return databasePath
}

// ScanFiles Traverse files on basepath and scraps and parse work data and saves the work on the DB
func ScanFiles(db *sql.DB, basepath string) {
	absPath, err := filepath.Abs(basepath)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !info.IsDir() && HasRJCode(info.Name()) {
			fmt.Printf("Scan file: %q\n", path)
			rjCode := GetRJCode(info.Name())
			work, err := fetcher.FetchWork(rjCode)
			if err != nil {
				return err
			}
			err = database.SaveWork(db, *work, path)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path: %v\n", err)
	}
}

// GetRJCode returns rjcode
func GetRJCode(filename string) string {
	foundRj := reCode.FindString(filename)

	return foundRj[2:]
}

// HasRJCode return true if has RJCode
func HasRJCode(filename string) bool {
	matched := reCode.MatchString(filename)
	return matched
}

// CreateSymlink creates a symlink of the file
func CreateSymlink(path string, newName string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Symlink(absPath, newName)

	if err != nil {
		log.Fatal(err)
	}
}

// FileExists checks if file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal(err)
	}
	return !info.IsDir()
}
