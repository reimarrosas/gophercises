package quizgame

import (
	"strings"
	"testing"
)

func TestConvertToQuiz(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		input := [][]string{{"2+2", "4"}}

		got, err := ConvertToQuiz(input)
		if err != nil {
			t.Errorf("expected no errors, got %v", err)
		}

		ql := got[0]
		gq, ga := ql.Question, ql.Answer
		wq, wa := input[0][0], input[0][1]

		if gq != wq || ga != wa {
			t.Errorf("expected Q: %s, expected A: %s, got Q: %s, got A: %s", wq, wa, gq, ga)
		}
	})

	t.Run("return error on invalid csv line", func(t *testing.T) {
		input := [][]string{{"2+2"}}

		_, err := ConvertToQuiz(input)
		if err == nil {
			t.Errorf("expected %v, got nil", ErrMalformedQuizCSV)
		}
	})
}

func TestAskQuiz(t *testing.T) {
	t.Run("successfully answer the question", func(t *testing.T) {
		input := [][]string{{"2+2", "4"}}
		q, _ := ConvertToQuiz(input)
		var sb strings.Builder
		r := strings.NewReader("4\n")

		got := askQuiz(q, r, &sb)
		var want uint = 1

		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	})

	t.Run("unsuccessfully answer the question", func(t *testing.T) {
		input := [][]string{{"2+2", "4"}}
		q, _ := ConvertToQuiz(input)
		var sb strings.Builder
		r := strings.NewReader("5\n")

		got := askQuiz(q, r, &sb)
		var want uint = 0

		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	})
}
