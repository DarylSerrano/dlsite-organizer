package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var dbDir string

var cmdRefresh = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh db",
	Long:  "Refresh database with information from dir arg or current dir",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
		var basePath string
		var err error
		if len(args) < 1 {
			basePath, err = os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			basePath = args[0]
		}

		fmt.Println("BasePath: " + basePath)
		fmt.Println("Databasepath: " + dbDir)
	},
}

var rootCmd = &cobra.Command{
	Use: "organizer",
}

func init() {
	basePath, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	cmdRefresh.PersistentFlags().StringVar(&dbDir, "db", basePath, "Dir where database is, default path is cwd/data.db")
	rootCmd.AddCommand(cmdRefresh)
}

func Execute() error {
	return rootCmd.Execute()
}
