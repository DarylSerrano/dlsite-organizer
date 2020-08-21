package main

import (
	"log"
	"os"
	"path/filepath"
)

func initialize(path string) {
	var dataPath = filepath.Join(path, "data")
	err := os.MkdirAll(dataPath, 0755)
	if err != nil {
		log.Panic(err)
	}
	var databasePath = filepath.Join(dataPath, "data.db")
	_, err = os.Create(databasePath)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	// if work, err := fetchWork("293003"); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(*work)
	// }
	// access()
	// showCurrentFolderFiles()
	// testCreateDir()

	// getRJCode()
	// hasRJCode()
	// createSymbolicLink()
	initialize(".")
}
