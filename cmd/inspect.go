package cmd

import (
	"fmt"
	"os"

	"github.com/bronzdoc/pak/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect package metadata",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Println("inspect subcommand needs and argument <artifact>")
			os.Exit(1)
		}

		artifactName := args[0]

		var label string
		if len(args) > 1 {
			label = args[1]
		}

		options := map[string]string{
			"label": label,
		}

		metadata, err := api.Inspect(artifactName, options)
		if err != nil {
			fmt.Println(errors.Wrap(err, "could not inspect metadata"))
			os.Exit(1)
		}

		fmt.Println(metadata)
	},
}

func init() {
	RootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
