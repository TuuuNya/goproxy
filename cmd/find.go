package cmd

import (
	"goproxy/internal/engine"

	"github.com/spf13/cobra"
)

var FinderArgs engine.FinderArgs

var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find proxy servers",
	Run: func(cmd *cobra.Command, args []string) {

		// Find proxy servers
		finder := engine.NewFinder(FinderArgs)
		finder.Find()
	},
}

func init() {
	FindCmd.Flags().StringSliceVarP(&FinderArgs.Type, "type", "t", []string{"http", "https", "socks5"}, "Proxy server types")
}
