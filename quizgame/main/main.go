package main

import (
	"flag"
	"fmt"
	"gophercises/quizgame"
	"os"
	"time"
)

func main() {
	var csvpath string
	flag.StringVar(&csvpath, "csv", "problems.csv", "Path to CSV Quiz")

	var duration string
	flag.StringVar(&duration, "duration", "30s", "Duration of the Quiz")

	flag.Parse()

	d, err := time.ParseDuration(duration)
	if err != nil {
		fatal(err)
	}

	rqd, err := quizgame.ReadCSV(csvpath)
	if err != nil {
		fatal(err)
	}

	q, err := quizgame.ConvertToQuiz(rqd)
	if err != nil {
		fatal(err)
	}

    res := "Quiz done!"
	score, err := quizgame.AskQuiz(q, d)
    if err != nil {
        res = "\nTimer ran out!"
    }
	fmt.Printf("%s Score: %d/%d\n", res, score, len(q))
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
