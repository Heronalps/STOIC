package cmd

import (
	"github.com/heronalps/STOIC/server"
	"github.com/spf13/cobra"
)

var (
	root        string
	registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register all images",
		Long:  `Register all image file names recursively from the root dir`,
		Run: func(cmd *cobra.Command, args []string) {
			server.RegisterImages(root)
		},
	}
)

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&root, "root", "r", "/opt", "The root dir of images")
}
