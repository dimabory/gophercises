package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dimabory/gophercises/3-adventure"
	"os"
	"strconv"
	"strings"
)

var (
	filename string
)

func init() {
	flag.StringVar(&filename, "story", "gopher.json", "story path")
	flag.Parse()
}

func main() {

	file, err := adventure.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	story, err := adventure.JsonStory(file)

	if err != nil {
		panic(err)
	}

	currentStory := story["intro"]

	askChan := make(chan string, len(story))
	quitChan := make(chan int, 1)

	ask(currentStory, askChan, quitChan)

	for {
		select {
		case c := <-askChan:
			if currentStory, ok := story[c]; ok {
				ask(currentStory, askChan, quitChan)
			}
		case <-quitChan:
			return
		}
	}

}

func ask(c adventure.Chapter, ch chan string, quit chan int) {

	printStory(c)

	if len(c.Options) == 0 {
		close(ch)
		quit <- 1
		return
	}

	printOptions(c.Options...)

	selectedStep := waitForAnswer(1, len(c.Options))

	var chapter string
	for i, v := range c.Options {
		if selectedStep-1 == i {
			chapter = v.Chapter
		}
	}

	ch <- chapter
}

func printStory(c adventure.Chapter) {
	fmt.Println(fmt.Sprintf(`
=== %s ===
%s
`,
		c.Title,
		strings.Join(c.Paragraphs, "\n"),
	))
}

func printOptions(options ...adventure.Option) {
	for i, v := range options {
		fmt.Printf("%d. %s\n", i+1, v.Text)
	}
}

func waitForAnswer(min, max int) int {
	fmt.Print("Please select your option.\n")

	input, _ := bufio.
		NewReader(os.Stdin).
		ReadString('\n')

	i, err := strconv.Atoi(strings.TrimSuffix(input, "\n"))
	if err != nil || i < min || i > max {
		fmt.Print("You've selected incorrect value.\n")
		return waitForAnswer(min, max)
	}

	return i
}
