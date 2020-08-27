package main

import (
	"log"

	"github.com/DarylSerrano/dlsite-organizer/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
