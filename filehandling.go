package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var reCode = regexp.MustCompile(`R(J|G)\d+`)

func showCurrentFolderFiles() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		log.Printf("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path: %v\n", err)
	}
}

func scanFiles(basepath string, db *sql.DB) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		log.Printf("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		log.Fatalf("error walking the path: %v\n", err)
	}
}

func testCreateDir() {
	err := os.Mkdir("testDir2", 0755)
	if err != nil {
		log.Fatal(err)
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

func createSymbolicLink() {
	absPath, err := filepath.Abs("foo")
	if err != nil {
		log.Panic(err)
	}

	err = os.Symlink(absPath, "foolink")

	if err != nil {
		log.Panic(err)
	}

	err = os.Rename("foolink", "links/foolink")
	if err != nil {
		log.Panic(err)
	}
}
