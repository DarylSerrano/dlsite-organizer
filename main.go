package main

import (
	"fmt"
	"os"

	"github.com/DarylSerrano/dlsite-organizer/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
