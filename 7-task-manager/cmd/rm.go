package cmd

import (
	"fmt"
	"github.com/dimabory/gophercises/7-task-manager/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove task from the list.",
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
			t := tasks[id-1]
			if ok := rmTask(t); ok == false {
				continue
			}
			fmt.Printf("You have deleted the \"%s\" task.\n", t.Value)
		}
	},
}

func rmTask(task db.Task) bool {
	err := db.Delete(task.Key)
	if err != nil {
		fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", task.Key, err)
		return false
	}
	return true
}
