package service

import (
	"fmt"
	"math-quiz/internal/config"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound  = fmt.Errorf("session not found")
	ErrQuestionNotFound = fmt.Errorf("question not found")
)

type QuizService struct {
	LenQuestions int
	sessionLock  sync.RWMutex
	sessions     map[string]*Session
}

func NewQuizService(cfg config.Config) *QuizService {
	qs := &QuizService{
		LenQuestions: cfg.LenQuestions,
		sessions:     make(map[string]*Session),
	}
	go qs.StartGC()

	return qs
}

func (qs *QuizService) StartGC() {
	for {
		<-time.After(time.Minute)
		qs.sessionLock.Lock()
		for k, v := range qs.sessions {
			if v.ExpiresAt.Before(time.Now()) {
				delete(qs.sessions, k)
			}
		}
		qs.sessionLock.Unlock()
	}
}

func (qs *QuizService) GetSession(sessionID string) *Session {
	qs.sessionLock.RLock()
	session := qs.sessions[sessionID]
	defer qs.sessionLock.RUnlock()
	return session
}

func (qs *QuizService) SessionExists(sessionID string) bool {
	qs.sessionLock.RLock()
	_, exists := qs.sessions[sessionID]
	defer qs.sessionLock.RUnlock()
	return exists
}

func (qs *QuizService) NewSessionID() string {
	return uuid.New().String()
}

func (qs *QuizService) GenerateAddSubQuestions(sessionID string) {
	session := NewSession()
	for range qs.LenQuestions {
		question := NewAddSubQuestion(4)
		session.Questions = append(session.Questions, question)
	}
	qs.AddSession(sessionID, session)
}

func (qs *QuizService) GenerateMultDivQuestions(sessionID string) {
	session := NewSession()
	for range qs.LenQuestions {
		question := NewMultDivQuestion(4)
		session.Questions = append(session.Questions, question)
	}
	qs.AddSession(sessionID, session)
}

func (qs *QuizService) AddSession(sessionID string, session *Session) {
	qs.sessionLock.Lock()
	session.ExpiresAt = time.Now().Add(time.Hour)
	qs.sessions[sessionID] = session
	defer qs.sessionLock.Unlock()
}
