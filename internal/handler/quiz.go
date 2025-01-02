package handler

import (
	"fmt"
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

func (h quizHandler) handleStartAddSub(w http.ResponseWriter, r *http.Request) error {

	sessionID := h.sessionIDFromRequest(r)
	if !h.service.SessionExists(sessionID) {
		sessionID = h.service.NewSessionID()
	}
	h.service.GenerateAddSubQuestions(sessionID)

	newCookie := &http.Cookie{
		Name:  "session-id",
		Value: sessionID,
	}
	http.SetCookie(w, newCookie)

	url := fmt.Sprintf("%s/quiz/0", h.BaseURL)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (h quizHandler) handleStartMultDiv(w http.ResponseWriter, r *http.Request) error {

	sessionID := h.sessionIDFromRequest(r)
	if !h.service.SessionExists(sessionID) {
		sessionID = h.service.NewSessionID()
	}
	h.service.GenerateMultDivQuestions(sessionID)

	newCookie := &http.Cookie{
		Name:  "session-id",
		Value: sessionID,
	}
	http.SetCookie(w, newCookie)

	url := fmt.Sprintf("%s/quiz/0", h.BaseURL)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (h quizHandler) handleQuiz(w http.ResponseWriter, r *http.Request) error {

	session := h.sessionFromRequest(r)
	if session == nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	idStr := chi.URLParam(r, "question_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	question, err := session.GetQuestionByID(id)
	if err != nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}
	questionLen := session.GetQuestionLen()

	props := quiz.QuestionProps{ID: id, QuestionLen: questionLen, Question: question, BaseURL: h.BaseURL}
	return quiz.Question(props).Render(r.Context(), w)
}

func (h quizHandler) handleAnswer(w http.ResponseWriter, r *http.Request) error {

	session := h.sessionFromRequest(r)
	if session == nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	idStr := chi.URLParam(r, "question_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	question, err := session.GetQuestionByID(id)
	if err != nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	answerIDStr := chi.URLParam(r, "answer_id")
	answerID, err := strconv.Atoi(answerIDStr)
	if err != nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	err = question.SetAnswer(answerID)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, fmt.Sprintf("%s/quiz/%d", h.BaseURL, id), http.StatusTemporaryRedirect)
	return nil
}

func (h quizHandler) handleResults(w http.ResponseWriter, r *http.Request) error {

	session := h.sessionFromRequest(r)
	if session == nil {
		url := fmt.Sprintf("%s/", h.BaseURL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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
