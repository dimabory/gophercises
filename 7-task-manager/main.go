package main

import (
	"github.com/dimabory/gophercises/7-task-manager/cmd"
	"github.com/dimabory/gophercises/7-task-manager/db"
	"github.com/mitchellh/go-homedir"
	"log"
	"path/filepath"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")

	if err := db.Init(dbPath); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
