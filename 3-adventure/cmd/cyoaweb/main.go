package main

import (
	"flag"
	"fmt"
	"github.com/dimabory/gophercises/3-adventure"
	"log"
	"net/http"
)

var (
	filename string
	port     int
)

func init() {
	flag.StringVar(&filename, "file", "gopher.json", "the JSON file with the CYOA story")
	flag.IntVar(&port, "port", 3000, "the port to start the CYOA web application on")

	flag.Parse()
}

func main() {
	fmt.Printf("Using the story in %s\n", filename)

	file, err := adventure.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	story, err := adventure.JsonStory(file)

	if err != nil {
		panic(err)
	}

	h := adventure.NewHandler(story)

	fmt.Println(fmt.Sprintf("Starting the server on port: %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))
}
