package main

import (
	"log"

	"github.com/spf13/cobra"
)

var pubsubDriver string
var numWorkers int
var redisAddr string
var redisDb int
var numMessage int
var debug bool

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

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List SQS queues",
	Args:  cobra.ExactArgs(0),
	Run:   runList,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&pubsubDriver, "driver", "gochannel", "Pub/Sub driver of: gochannel, redis, sqs")
	rootCmd.PersistentFlags().IntVar(&numWorkers, "workers", 1, "Number of workers")
	rootCmd.PersistentFlags().StringVar(&redisAddr, "redisaddr", "localhost:6379", "Redis address")
	rootCmd.PersistentFlags().IntVar(&redisDb, "redisdb", 0, "Redis db")
	rootCmd.PersistentFlags().IntVar(&numMessage, "messages", 10, "Number of messages to publish")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable/Disable debug log message")
	rootCmd.AddCommand(cmdRouter)
	rootCmd.AddCommand(cmdPubSub)
	rootCmd.AddCommand(cmdList)
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
