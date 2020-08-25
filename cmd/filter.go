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

var isSfw bool

var cmdRootFilter = &cobra.Command{
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

var cmdSfwFilter = &cobra.Command{
	Use:   "sfw",
	Short: "Filter work by sfw/nsfw",
	Long:  "Filter work by sfw/nsfw",
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

		filter.FilterBySfw(db, isSfw, *basePath)
	},
}

func init() {
	cmdSfwFilter.Flags().BoolVarP(&isSfw, "isSFW", "s", true, "Shoild filter by SFW Works")
	cmdRootFilter.AddCommand(cmdSfwFilter)
	rootCmd.AddCommand(cmdRootFilter)
}
