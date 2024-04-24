package main

import (
	"log"

	"github.com/alitto/pond"
	"github.com/spf13/cobra"
)

var pubsubDriver string
var numWorkers int
var redisAddr string
var redisDb int
var numMessage int

var emailWorker *pond.WorkerPool
var imosWorker *pond.WorkerPool

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
	rootCmd.PersistentFlags().StringVar(&pubsubDriver, "driver", "gochannel", "Pub/Sub driver of: gochannel, redis")
	rootCmd.PersistentFlags().IntVar(&numWorkers, "workers", 2, "Number of workers")
	rootCmd.PersistentFlags().StringVar(&redisAddr, "redisaddr", "localhost:6379", "Redis address")
	rootCmd.PersistentFlags().IntVar(&redisDb, "redisdb", 0, "Redis db")
	rootCmd.PersistentFlags().IntVar(&numMessage, "messages", 10, "Number of messages to publish")
	rootCmd.AddCommand(cmdRouter)
	rootCmd.AddCommand(cmdPubSub)
}

func main() {
	// create an unbuffered (blocking) pool with a number of workers
	emailWorker = pond.New(numWorkers, 0, pond.MinWorkers(numWorkers))
	defer emailWorker.StopAndWait()
	imosWorker = pond.New(numWorkers, 0, pond.MinWorkers(numWorkers))
	defer imosWorker.StopAndWait()
	// run root command
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("root command error")
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
