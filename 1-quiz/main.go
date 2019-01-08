package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 10, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open a CSV file: %s", *csvFilename), 1)
	}

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.", 1)
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	// problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		answerChan := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChan <- strings.ToLower(strings.TrimSpace(answer))
		}()

		select {
		case <-timer.C:
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

	return shuffle(result)
}

func shuffle(src []problem) (final []problem) {
	final = make([]problem, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}

	return
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
