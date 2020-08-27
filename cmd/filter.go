package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/DarylSerrano/dlsite-organizer/internal/database"
	"github.com/DarylSerrano/dlsite-organizer/internal/filehandler"
	"github.com/DarylSerrano/dlsite-organizer/internal/filter"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

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
			log.Fatal(err)
		}

		filter.BySfw(db, true, *basePath)
	},
}

var cmdSfwFilter = &cobra.Command{
	Use:   "sfw",
	Short: "Filter work by sfw",
	Long:  "Filter work by sfw",
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
			log.Fatal(err)
		}

		filter.BySfw(db, true, *basePath)
	},
}

var cmdNsfwFilter = &cobra.Command{
	Use:   "nsfw",
	Short: "Filter work by Nsfw",
	Long:  "Filter work by Nsfw",
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
			log.Fatal(err)
		}

		filter.BySfw(db, false, *basePath)
	},
}

var cmdTagFilter = &cobra.Command{
	Use:   "tags",
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
			log.Fatal(err)
		}

		tags := filter.GetAllTags(db)

		if all {
			for _, tag := range tags {
				filter.ByTag(db, tag, *basePath)
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
				Label:     "What Tag",
				Items:     tags,
				Templates: templates,
				Size:      5,
				Searcher:  searcher,
			}

			i, _, err := prompt.Run()

			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}

			filter.ByTag(db, tags[i], *basePath)
		}
	},
}

var cmdCircleFilter = &cobra.Command{
	Use:   "circles",
	Short: "Filter by circles",
	Long:  "FIlter by circles saved on the db",
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
			log.Fatal(err)
		}

		circles := filter.GetAllCircles(db)

		if all {
			for _, circle := range circles {
				filter.ByCircle(db, circle, *basePath)
			}
		} else {
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U000027A1 {{ .Name | cyan }}",
				Inactive: "  {{ .Name | cyan }}",
			}

			searcher := func(input string, index int) bool {
				tag := circles[index]
				name := strings.Replace(strings.ToLower(tag.Name), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			}

			prompt := promptui.Select{
				Label:     "What Circle",
				Items:     circles,
				Templates: templates,
				Size:      5,
				Searcher:  searcher,
			}

			i, _, err := prompt.Run()

			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}

			filter.ByCircle(db, circles[i], *basePath)
		}
	},
}

var cmdVAFilter = &cobra.Command{
	Use:   "vas",
	Short: "Filter by Voice actors",
	Long:  "FIlter by Voice actors saved on the db",
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
			log.Fatal(err)
		}

		voiceActos := filter.GetAllVoiceActors(db)

		if all {
			for _, voiceActor := range voiceActos {
				filter.ByVoiceActor(db, voiceActor, *basePath)
			}
		} else {
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U000027A1 {{ .Name | cyan }}",
				Inactive: "  {{ .Name | cyan }}",
			}

			searcher := func(input string, index int) bool {
				tag := voiceActos[index]
				name := strings.Replace(strings.ToLower(tag.Name), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			}

			prompt := promptui.Select{
				Label:     "What Voice Actor",
				Items:     voiceActos,
				Templates: templates,
				Size:      5,
				Searcher:  searcher,
			}

			i, _, err := prompt.Run()

			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}

			filter.ByVoiceActor(db, voiceActos[i], *basePath)
		}
	},
}

func configureFilter() {
	cmdTagFilter.Flags().BoolVarP(&all, "all", "a", false, "Filter all instead of select one")
	cmdCircleFilter.Flags().BoolVarP(&all, "all", "a", false, "Filter all instead of select one")
	cmdVAFilter.Flags().BoolVarP(&all, "all", "a", false, "Filter all instead of select one")
	cmdRootFilter.AddCommand(cmdSfwFilter, cmdNsfwFilter, cmdTagFilter, cmdCircleFilter, cmdVAFilter)
	rootCmd.AddCommand(cmdRootFilter)
}
