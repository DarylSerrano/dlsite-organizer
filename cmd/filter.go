package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/DarylSerrano/dlsite-organizer/internal/database"
	"github.com/DarylSerrano/dlsite-organizer/internal/filehandler"
	"github.com/DarylSerrano/dlsite-organizer/internal/filter"
	"github.com/spf13/cobra"
)

var cmdFilter = &cobra.Command{
	Use:   "filter",
	Short: "Filter works, default filter by sfw",
	Long:  "Filter works by circle, va, tags..., the default is to filter by sfw",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		basePath, err := getBasePath(args)
		if err != nil {
			log.Fatal(err)
		}
		databasePath := filehandler.CreateDBFile(dbDir)

		fmt.Println("BasePath: " + *basePath)
		fmt.Println("Databasepath: " + dbDir)

		db, err := database.OpenDB(databasePath)
		defer db.Close()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		filter.FilterBySfw(db, true, *basePath)
	},
}

func init() {
	rootCmd.AddCommand(cmdFilter)
}
