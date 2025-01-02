package service

import "time"

type Session struct {
	Questions []*QuizQuestion
	ExpiresAt time.Time
}

type QuizResult struct {
	Grade            int
	CorrectAnswers   int
	IncorrectAnswers int
}

func NewSession() *Session {
	return &Session{
		Questions: make([]*QuizQuestion, 0),
	}
}

func (s *Session) GetQuestionByID(id int) (*QuizQuestion, error) {
	if id < 0 || id >= len(s.Questions) {
		return nil, ErrQuestionNotFound
	}
	question := s.Questions[id]
	return question, nil
}

func (s *Session) GetQuestionLen() int {
	return len(s.Questions)
}

func (s *Session) GenerateResult() *QuizResult {
	baseGrade := 10
	result := &QuizResult{}

	for _, q := range s.Questions {
		if q.IsCorrect {
			result.CorrectAnswers++
		} else {
			result.IncorrectAnswers++
		}
	}

	result.Grade = (result.CorrectAnswers * baseGrade) / s.GetQuestionLen()

	return result
}
