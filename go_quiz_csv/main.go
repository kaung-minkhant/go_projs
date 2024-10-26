package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

func parseProblems(lines [][]string) []Problem {
	// 1. go over lines and parse the lines based on problem struct
	var problems []Problem
	for _, line := range lines {
		problems = append(problems, Problem{
			Question: line[0],
			Answer:   line[1],
		})
	}
	return problems
}

func pullProblems(filename string) ([]Problem, error) {
	// 1. open the file
	file, err := os.Open(filename)
  defer file.Close()
	// 2. handle err of opening file
	if err != nil {
    return nil, fmt.Errorf("Opening quiz file error: %s\n", err)
	}
	// 3. use csv reader to read lines of records
	csvR := csv.NewReader(file)
	records, err := csvR.ReadAll()
	if err != nil {
    return nil, fmt.Errorf("Reading csv file error: %s\n", err)
	}
	// 4. Parse the csv file (parse csv)
	return parseProblems(records), nil
}

func main() {
	// 1. get the name of the problem file
	quiz := flag.String("f", "quiz.csv", "path of the csv quiz file")
	// 2. set the timer value
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()
	// 3. pull the problems from the file (problem puller)
	problems, err := pullProblems(*quiz)
	// 4. handle error from pull
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 5. initial number of correct ansers
	correctAnswers := 0
	// 6. create timer for quizes
	quizTimer := time.NewTimer(time.Duration(*timer) * time.Second)
	// 7. Channel to receive answers
	ansC := make(chan string)
	// 8. loop through the problems, accept answers, reset timer if necesary and timeout
  fmt.Printf("You have %d seconds to answer all questions.\n", *timer)
  problemLoop:
	for i, problem := range problems {
    fmt.Printf("Problem %d: %s = ", i+1, problem.Question)

    go func() {
      var answer string
      fmt.Scanf("%s", &answer)
      ansC <- answer
    }()

    select {
    case <-quizTimer.C:
      fmt.Println()
      // close(ansC)
      break problemLoop
    case answer := <- ansC:
      if answer == problem.Answer {
        correctAnswers++
      }
      if i == len(problems)-1 {
        close(ansC)
      }
    }
	}
	// 9. calcuate the result and print
  fmt.Printf("Your result is %d out of %d.\n", correctAnswers, len(problems))

  fmt.Printf("Press Enter to exist")
  <-ansC
}
