package cmd

import (
	"fmt"
	"github.com/dimabory/gophercises/7-task-manager/db"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new task to the list.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("No task provided.")
			return
		}

		task := strings.Join(args, " ")

		_, err := db.Add(db.Task{Value: task})
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}
