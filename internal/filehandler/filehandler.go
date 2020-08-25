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

func CreateDBFile(path string) string {
	var databasePath = filepath.Join(path, "data.db")
	if !FileExists(databasePath) {
		log.Print("Database doesnt exists, creating", databasePath)
		_, err := os.Create(databasePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return databasePath
}

func ScanFiles(db *sql.DB, basepath string) {
	absPath, err := filepath.Abs(basepath)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("abs path: ", absPath)

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		log.Printf("visited file or dir: %q\n", path)

		if info.IsDir() {
			return nil
		}

		if !info.IsDir() && HasRJCode(info.Name()) {
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

func GetRJCode(filename string) string {
	foundRj := reCode.FindString(filename)

	return foundRj[2:]
}

func HasRJCode(filename string) bool {
	matched := reCode.MatchString(filename)
	return matched
}

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
