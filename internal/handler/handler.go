package handler

import (
	"io/fs"
	"log/slog"
	"math-quiz/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

type Dependencies struct {
	QuizService *service.QuizService
	AssetsFS    fs.FS
	BaseURL     string
}

func RegisterRoutes(r *chi.Mux, deps Dependencies) {
	home := homeHandler{BaseURL: deps.BaseURL}
	quiz := quizHandler{service: deps.QuizService, BaseURL: deps.BaseURL}

	r.Get("/", handler(home.handleIndex))

	r.Get("/startaddsub", handler(quiz.handleStartAddSub))
	r.Get("/startmultdiv", handler(quiz.handleStartMultDiv))
	r.Get("/quiz/{question_id}", handler(quiz.handleQuiz))
	r.Get("/quiz/{question_id}/{answer_id}", handler(quiz.handleAnswer))
	r.Get("/results", handler(quiz.handleResults))

	r.Handle("/web/public/assets/*", http.FileServerFS(deps.AssetsFS))
}

func handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("Error handling request", "err", err)
		}
	}
}
