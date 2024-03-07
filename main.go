package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "wk-api server",
	Run:   runRoot,
}

var cmdPubSub = &cobra.Command{
	Use:   "pubsub",
	Short: "Test watermill with low level pubsub functions",
	Args:  cobra.ExactArgs(0),
	Run:   runPubSub,
}

var cmdRouter = &cobra.Command{
	Use:   "router",
	Short: "Test watermill with router functions",
	Args:  cobra.ExactArgs(0),
	Run:   runRouter,
}

func init() {
	rootCmd.AddCommand(cmdRouter)
	rootCmd.AddCommand(cmdPubSub)
}

func main() {
	// run root command
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("root command error")
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
