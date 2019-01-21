package cmd

import (
	"fmt"
	"github.com/dimabory/gophercises/7-task-manager/db"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(completedCmd)
}

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "Show all completed tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.GetCompleted()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no completed tasks!")
			return
		}

		fmt.Println("You have completed the following tasks:")
		for _, t := range tasks {
			fmt.Println(t)
		}

	},
}
