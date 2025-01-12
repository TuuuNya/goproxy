package cmd

import (
	"goproxy/internal/engine"
	"time"

	"github.com/spf13/cobra"
)

var ServerArgs engine.ServerArgs

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start proxy server",
	Run: func(cmd *cobra.Command, args []string) {
		duration, _ := cmd.Flags().GetDuration("max_delay")
		ServerArgs.MaxDelay = duration
		server := engine.NewServer(ServerArgs)
		server.Serve()
	},
}

func init() {
	ServeCmd.Flags().StringVarP(&ServerArgs.Type, "type", "t", "socks5", "Proxy server type")
	ServeCmd.Flags().IntVarP(&ServerArgs.Port, "port", "p", 7777, "Port to listen on")
	ServeCmd.Flags().Duration("max_delay", 10*time.Second, "Maximum delay for proxy server")
	ServeCmd.Flags().IntVarP(&ServerArgs.CheckPoolSize, "check_pool_size", "c", 100000, "Size of the pool to check proxies")
}
