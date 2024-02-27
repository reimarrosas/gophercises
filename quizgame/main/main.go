package main

import (
	"flag"
	"fmt"
	"gophercises/quizgame"
	"os"
)

func main() {
    var csvpath string
    flag.StringVar(&csvpath, "csv", "problems.csv", "Path to CSV Quiz")
    flag.Parse()

    rqd, err := quizgame.ReadCSV(csvpath)
    if err != nil {
        fatal(err)
    }

    q, err := quizgame.ConvertToQuiz(rqd)
    if err != nil {
        fatal(err)
    }

    score := quizgame.AskQuiz(q)
	fmt.Printf("Quiz done! Score: %d/%d\n", score, len(q))
}

func fatal(err error) {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
