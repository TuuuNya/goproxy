package main

import (
	"goproxy/cmd"
	"goproxy/pkg/logger"

	"github.com/spf13/cobra"
)

type cmdArgs struct {
	Debug bool
}

func main() {
	log := logger.GetLogger()

	var rootArgs cmdArgs

	rootCmd := &cobra.Command{
		Use:   "goproxy",
		Short: "goproxy is a Go language implementation of a proxy broker.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetLogLevel(rootArgs.Debug)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&rootArgs.Debug, "debug", "d", false, "Enable debug mode")

	rootCmd.AddCommand(cmd.FindCmd)
	rootCmd.AddCommand(cmd.ServeCmd)

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Failed to execute root command")
	}
}
