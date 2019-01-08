package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
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

	correct := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == problem.answer {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) (result []problem) {
	result = make([]problem, len(lines))

	for i, line := range lines {
		result[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
