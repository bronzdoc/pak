package cmd

import (
	"fmt"
	"os"

	"github.com/bronzdoc/pak/api"
	"github.com/spf13/cobra"
)

// promoteCmd represents the promote command
var promoteCmd = &cobra.Command{
	Use:   "promote <ARTIFACT> <LABEL>",
	Short: "Promote an artifact",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Println("promote subcommand needs and argument <artifact>")
			os.Exit(1)
		}

		artifactName := args[0]

		var label string
		if len(args) > 1 {
			label = args[1]
		}

		options := map[string]interface{}{
			"label": label,
		}

		artifactName, err := api.Promote(artifactName, options)
		if err != nil {
			fmt.Printf("failed to promote %s: %s", artifactName, err)
			os.Exit(1)
		}

		fmt.Println(artifactName)
	},
}

func init() {
	RootCmd.AddCommand(promoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promoteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
