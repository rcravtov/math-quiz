package service

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
)

var (
	ErrAnswerNotFound = fmt.Errorf("answer not found")
)

type QuizQuestion struct {
	Question         string
	Answers          []string
	IsAnswered       bool
	IsCorrect        bool
	CorrectAnswerID  int
	SelectedAnswerID int
}

func (q *QuizQuestion) SetAnswer(answerID int) error {

	if answerID < 0 || answerID >= len(q.Answers) {
		return ErrAnswerNotFound
	}

	q.IsAnswered = true
	q.SelectedAnswerID = answerID
	q.IsCorrect = q.CorrectAnswerID == q.SelectedAnswerID

	return nil
}

func NewAddSubQuestion(lenAnswers int, max int) *QuizQuestion {
	var (
		left   = 0
		right  = 0
		result = 0
	)

	signs := []string{"+", "-"}
	sign := signs[rand.Intn(len(signs))]

	for {
		left = rand.Intn(max) + 1
		right = rand.Intn(max) + 1

		if sign == "+" {
			result = left + right
		} else {
			result = left - right
		}

		if result >= 0 && result <= max {
			break
		}
	}

	answers, correctAnswerID := generateAnswers(lenAnswers, 100, result)

	q := &QuizQuestion{
		Question:        fmt.Sprintf("%d %s %d", left, sign, right),
		Answers:         answers,
		CorrectAnswerID: correctAnswerID,
	}

	return q
}

func NewMultDivQuestion(lenAnswers int) *QuizQuestion {
	var (
		left   = 0
		right  = 0
		result = 0
	)

	signs := []string{"*", ":"}
	sign := signs[rand.Intn(len(signs))]

	left = rand.Intn(8) + 2
	right = rand.Intn(8) + 2

	if sign == "*" {
		result = left * right
	} else {
		result, left = left, left*right
	}

	answers, correctAnswerID := generateAnswers(lenAnswers, result+10, result)

	q := &QuizQuestion{
		Question:        fmt.Sprintf("%d %s %d", left, sign, right),
		Answers:         answers,
		CorrectAnswerID: correctAnswerID,
	}

	return q
}

func generateAnswers(lenAnswers int, max int, correctAnswer int) ([]string, int) {
	if lenAnswers <= 0 {
		lenAnswers = 2
	}

	answers := make([]int, 0, lenAnswers)
	answers = append(answers, correctAnswer)

	for {
		answer := rand.Intn(max) + 1
		if slices.Contains(answers, answer) {
			continue
		}

		answers = append(answers, answer)
		if len(answers) == lenAnswers {
			break
		}
	}
	swapIndex := rand.Intn(lenAnswers)
	answers[0], answers[swapIndex] = answers[swapIndex], answers[0]

	strAnswers := make([]string, lenAnswers)
	for i, v := range answers {
		strAnswers[i] = strconv.Itoa(v)
	}

	return strAnswers, swapIndex
}
