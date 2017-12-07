package cmd

import (
	"fmt"
	"os"

	"github.com/bronzdoc/pak/api"
	"github.com/bronzdoc/pak/pakfile"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an artifact",
	Run: func(cmd *cobra.Command, args []string) {
		pakfile, err := pakfile.Factory()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(api.Build(pakfile))

	},
}

func init() {
	RootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
