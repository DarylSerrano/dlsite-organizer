package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DarylSerrano/dlsite-organizer/internal/database"
	"github.com/DarylSerrano/dlsite-organizer/internal/filehandler"
	"github.com/DarylSerrano/dlsite-organizer/internal/filter"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var isSfw bool
var all bool

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

var cmdTagFilter = &cobra.Command{
	Use:   "tag",
	Short: "Filter by tag",
	Long:  "FIlter by tag saved on the db",
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

		tags := filter.GetAllTags(db)

		if all {
			for _, tag := range tags {
				filter.FilterByTag(db, tag, *basePath)
			}
		} else {
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U000027A1 {{ .Name | cyan }}",
				Inactive: "  {{ .Name | cyan }}",
			}

			searcher := func(input string, index int) bool {
				tag := tags[index]
				name := strings.Replace(strings.ToLower(tag.Name), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			}

			prompt := promptui.Select{
				Label:     "Select Tag",
				Items:     tags,
				Templates: templates,
				Size:      5,
				Searcher:  searcher,
			}

			i, _, err := prompt.Run()

			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}

			id := tags[i].ID
			log.Println("You picked: index", i, "result: ", tags[i].Name, " tag ID: ", id)

			filter.FilterByTag(db, tags[i], *basePath)
		}

	},
}

func init() {
	cmdSfwFilter.Flags().BoolVarP(&isSfw, "isSFW", "s", true, "Should filter by SFW Works")
	cmdTagFilter.Flags().BoolVarP(&all, "all", "a", false, "Filter all instead of select one")
	cmdRootFilter.AddCommand(cmdSfwFilter, cmdTagFilter)
	rootCmd.AddCommand(cmdRootFilter)
}
