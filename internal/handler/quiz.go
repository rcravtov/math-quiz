package handler

import (
	"math-quiz/internal/service"
	"math-quiz/internal/view/quiz"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type quizHandler struct {
	service *service.QuizService
	BaseURL string
}

func (h quizHandler) handleStartAddSub10(w http.ResponseWriter, r *http.Request) error {

	sessionID := h.sessionIDFromRequest(r)
	if !h.service.SessionExists(sessionID) {
		sessionID = h.service.NewSessionID()
	}
	session := h.service.GenerateAddSubQuestions(sessionID, 10)

	newCookie := &http.Cookie{
		Name:  "session-id",
		Value: sessionID,
	}
	http.SetCookie(w, newCookie)

	props := quiz.QuestionProps{
		ID:          0,
		QuestionLen: session.GetQuestionLen(),
		Question:    session.Questions[0],
		BaseURL:     h.BaseURL,
	}
	return quiz.Question(props).Render(r.Context(), w)
}

func (h quizHandler) handleStartAddSub100(w http.ResponseWriter, r *http.Request) error {

	sessionID := h.sessionIDFromRequest(r)
	if !h.service.SessionExists(sessionID) {
		sessionID = h.service.NewSessionID()
	}
	session := h.service.GenerateAddSubQuestions(sessionID, 100)

	newCookie := &http.Cookie{
		Name:  "session-id",
		Value: sessionID,
	}
	http.SetCookie(w, newCookie)

	props := quiz.QuestionProps{
		ID:          0,
		QuestionLen: session.GetQuestionLen(),
		Question:    session.Questions[0],
		BaseURL:     h.BaseURL,
	}
	return quiz.Question(props).Render(r.Context(), w)
}

func (h quizHandler) handleStartMultDiv(w http.ResponseWriter, r *http.Request) error {

	sessionID := h.sessionIDFromRequest(r)
	if !h.service.SessionExists(sessionID) {
		sessionID = h.service.NewSessionID()
	}
	session := h.service.GenerateMultDivQuestions(sessionID)

	newCookie := &http.Cookie{
		Name:  "session-id",
		Value: sessionID,
	}
	http.SetCookie(w, newCookie)

	props := quiz.QuestionProps{
		ID:          0,
		QuestionLen: session.GetQuestionLen(),
		Question:    session.Questions[0],
		BaseURL:     h.BaseURL,
	}
	return quiz.Question(props).Render(r.Context(), w)
}

func (h quizHandler) handleQuiz(w http.ResponseWriter, r *http.Request) error {

	session := h.sessionFromRequest(r)
	if session == nil {
		return service.ErrSessionNotFound
	}

	idStr := chi.URLParam(r, "question_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	question, err := session.GetQuestionByID(id)
	if err != nil {
		return err
	}
	questionLen := session.GetQuestionLen()

	props := quiz.QuestionProps{
		ID:          id,
		QuestionLen: questionLen,
		Question:    question,
		BaseURL:     h.BaseURL,
	}
	return quiz.Question(props).Render(r.Context(), w)
}

func (h quizHandler) handleAnswer(w http.ResponseWriter, r *http.Request) error {

	session := h.sessionFromRequest(r)
	if session == nil {
		return service.ErrSessionNotFound
	}

	idStr := chi.URLParam(r, "question_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	question, err := session.GetQuestionByID(id)
	if err != nil {
		return err
	}

	answerIDStr := chi.URLParam(r, "answer_id")
	answerID, err := strconv.Atoi(answerIDStr)
	if err != nil {
		return err
	}

	err = question.SetAnswer(answerID)
	if err != nil {
		return err
	}

	// show results after last question
	if id == session.GetQuestionLen()-1 {
		return h.handleResults(w, r)
	}

	nextID := id + 1
	nextQuestion, err := session.GetQuestionByID(nextID)
	if err != nil {

	}
	props := quiz.QuestionProps{
		ID:          nextID,
		QuestionLen: session.GetQuestionLen(),
		Question:    nextQuestion,
		BaseURL:     h.BaseURL,
	}
	return quiz.Question(props).Render(r.Context(), w)
}

func (h quizHandler) handleResults(w http.ResponseWriter, r *http.Request) error {

	session := h.sessionFromRequest(r)
	if session == nil {
		return nil
	}

	result := session.GenerateResult()
	props := quiz.ResultProps{Result: result, Questions: session.Questions, BaseURL: h.BaseURL}
	return quiz.Results(props).Render(r.Context(), w)
}

func (h quizHandler) sessionIDFromRequest(r *http.Request) string {
	sessionID := (r.Context().Value(SessionKey("session-id"))).(string)
	return sessionID
}

func (h quizHandler) sessionFromRequest(r *http.Request) *service.Session {
	sessionID := h.sessionIDFromRequest(r)
	return h.service.GetSession(sessionID)
}
