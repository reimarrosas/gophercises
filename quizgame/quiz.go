package quizgame

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	ErrCannotOpenCSV    = errors.New("cannot open csv file")
	ErrCannotReadCSV    = errors.New("cannot read csv file")
	ErrMalformedQuizCSV = errors.New("malformed quiz csv")
	ErrTimerRanOut      = errors.New("timer ran out")
)

var m sync.Mutex

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

func askQuiz(q Quiz, d time.Duration, r io.Reader, w io.Writer) (uint, error) {
	s := bufio.NewScanner(r)
	var score uint

	a := make(chan string)
	defer close(a)

	t := time.NewTicker(d)
	defer t.Stop()

	for _, l := range q {
		go func(a chan string, s *bufio.Scanner, w io.Writer) {
			fmt.Fprintf(w, "%s: ", l.Question)

            m.Lock()
			for s.Scan() {
                a <- s.Text()
            }
            m.Unlock()
		}(a, s, w)

		select {
		case d := <-a:
			if d == l.Answer {
				score++
			}
		case <-t.C:
			return score, ErrTimerRanOut
		}
	}

	return score, nil
}

func AskQuiz(q Quiz, d time.Duration) (uint, error) {
	return askQuiz(q, d, os.Stdin, os.Stdout)
}
