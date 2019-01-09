package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	flagFilePath string
	flagRandom   bool
	flagTime     int
)

type problem struct {
	question string
	answer   string
}

func init() {
	// flagFilePath := flag.String("csv", "problems.csv", "path/to/csv_file")

	flag.StringVar(&flagFilePath, "csv", "problems.csv", "path/to/csv_file")
	flag.IntVar(&flagTime, "time", 10, "the time limit for the quiz in seconds")
	flag.BoolVar(&flagRandom, "shuffle", false, "shuffle questions?")

	flag.Parse()
}

func main() {

	lines := loadFile(flagFilePath)

	// block until user presses enter
	fmt.Print(fmt.Sprintf("Press [Enter] to start test. You will have %d seconds to finish the quiz.", flagTime))
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	problems := parseLines(lines)

	if flagRandom {
		problems = shuffleQuestions(problems)
	}

	timer := time.NewTimer(time.Duration(flagTime) * time.Second)
	correct := 0

	answerChan := make(chan string, len(problems))

	// problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChan <- strings.ToLower(strings.TrimSpace(answer))
		}()

		select {
		case <-timer.C:
			close(answerChan)
			fmt.Println()
			// break problemLoop
			return
		case answer := <-answerChan:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		result[i] = problem{
			question: line[0],
			answer:   strings.ToLower(strings.TrimSpace(line[1])),
		}
	}

	return result
}

func shuffleQuestions(src []problem) (final []problem) {
	final = make([]problem, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}

	return
}

func loadFile(filepath string) [][]string {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Couldn't open a CSV file: %s", filepath))
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Couldn't parse a CSV file: %s", filepath))
	}
	return records
}
