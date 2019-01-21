package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a CLI task manager.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ParseArgs(args ...string) (ids []int) {
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Failed to parse the argument:", arg)
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
