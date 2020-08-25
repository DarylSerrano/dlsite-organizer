package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var reCode = regexp.MustCompile(`R(J|G)\d+`)

func scanFiles(basepath string, db *sql.DB) {
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

		if !info.IsDir() && hasRJCode(info.Name()) {
			rjCode := getRJCode(info.Name())
			work, err := fetchWork(rjCode)
			if err != nil {
				return err
			}
			err = saveWork(db, *work, path)
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

func getRJCode(filename string) string {
	foundRj := reCode.FindString(filename)

	return foundRj[2:]
}

func hasRJCode(filename string) bool {
	matched := reCode.MatchString(filename)
	return matched
}

func createSymlink(path string, newName string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Symlink(absPath, newName)

	if err != nil {
		log.Fatal(err)
	}
}

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
