package cmd

import (
	"fmt"
	"github.com/dimabory/gophercises/7-task-manager/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Perform a task from the list.",
	Run: func(cmd *cobra.Command, args []string) {
		ids := ParseArgs(args...)

		tasks, err := db.GetAll()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}

		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number:", id)
				continue
			}
			if ok := markAsRead(tasks[id-1]); ok == false {
				continue
			}
			fmt.Printf("Marked \"%d\" as completed.\n", id)
		}
	},
}

func markAsRead(task db.Task) bool {
	if err := db.Complete(task); err != nil {
		fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", task.Key, err)
		return false
	}
	return true
}
