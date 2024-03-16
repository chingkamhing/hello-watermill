package main

import (
	"log"

	"github.com/alitto/pond"
	"github.com/spf13/cobra"
)

const numEmailWorkers = 2

var emailWorker *pond.WorkerPool

const numImosWorkers = 2

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
	rootCmd.AddCommand(cmdRouter)
	rootCmd.AddCommand(cmdPubSub)
}

func main() {
	// create an unbuffered (blocking) pool with a number of workers
	emailWorker = pond.New(numEmailWorkers, 0, pond.MinWorkers(numEmailWorkers))
	defer emailWorker.StopAndWait()
	imosWorker = pond.New(numImosWorkers, 0, pond.MinWorkers(numImosWorkers))
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
