package quizgame

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

// Read File
// Display Quiz line by line
// Prompt user answer by line
// Tally correct answers
// Display score

var ErrCannotOpenCSV = errors.New("cannot open csv file")
var ErrCannotReadCSV = errors.New("cannot read csv file")
var ErrMalformedQuizCSV = errors.New("malformed quiz csv")

type rawQuizData [][]string

type Quiz []quizLine

type quizLine struct {
	Question, Answer string
}

func ReadCSV(filename string) (rawQuizData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, ErrCannotOpenCSV
	}

	r := csv.NewReader(f)
	rqd, err := r.ReadAll()
	if err != nil {
		return nil, ErrCannotReadCSV
	}

	return rqd, nil
}

func ConvertToQuiz(rqd rawQuizData) (Quiz, error) {
	var res Quiz

	for _, l := range rqd {
        if len(l) != 2 {
            return nil, ErrMalformedQuizCSV
        }

		q, a := l[0], l[1]
		res = append(res, quizLine{q, a})
	}

	return res, nil
}

func askQuiz(q Quiz, r io.Reader, w io.Writer) {
	s := bufio.NewScanner(r)
	var score uint

	for _, l := range q {
		fmt.Fprintf(w, "%s: ", l.Question)

		s.Scan()
		a := s.Text()

		if a == l.Answer {
			score++
		}
	}

    fmt.Fprintf(w, "Quiz done! Score: %d/%d\n", score, len(q))
}

func AskQuiz(q Quiz) {
    askQuiz(q, os.Stdin, os.Stdout)
}
