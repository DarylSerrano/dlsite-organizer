package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/DarylSerrano/dlsite-organizer/internal/database"
	"github.com/DarylSerrano/dlsite-organizer/internal/filehandler"
	"github.com/spf13/cobra"
)

var dbDir string

var cmdRefresh = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh db",
	Long:  "Refresh database with information from dir arg or current dir",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		basePath, err := getBasePath(args)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("BasePath: " + *basePath)
		fmt.Println("Databasepath: " + dbDir)
		databasePath := filehandler.CreateDBFile(dbDir)
		db, err := database.OpenDB(databasePath)
		defer db.Close()
		if err != nil {
			log.Fatal(err)
		}
		filehandler.ScanFiles(db, *basePath)
	},
}

var rootCmd = &cobra.Command{
	Use: "dlsite-organizer",
}

func getBasePath(args []string) (*string, error) {
	var basePath string
	var err error
	if len(args) < 1 {
		basePath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
		return &basePath, nil
	}

	basePath = args[0]
	return &basePath, nil
}

func init() {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.PersistentFlags().StringVar(&dbDir, "db", basePath, "Dir where database is")
	rootCmd.AddCommand(cmdRefresh)
}

// Execute runs commandline
func Execute() error {
	return rootCmd.Execute()
}
