package app

import (
	"context"
	"flag"
	"log/slog"
	assets "math-quiz"
	"math-quiz/internal/config"
	"math-quiz/internal/handler"
	"math-quiz/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(ctx context.Context) error {

	var (
		addr       string
		baseURL    string
		nquestions int
	)

	flag.StringVar(&addr, "addr", ":3000", "server address")
	flag.StringVar(&baseURL, "baseurl", "", "base url")
	flag.IntVar(&nquestions, "nquestions", 20, "number of questions")
	flag.Parse()

	cfg := config.NewConfig(addr, baseURL, nquestions)
	srv := service.NewQuizService(cfg)

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(handler.SessionIDToContext)
	handler.RegisterRoutes(r, handler.Dependencies{
		AssetsFS:    assets.AssetsFS,
		QuizService: srv,
		BaseURL:     baseURL,
	})

	s := http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down server")
		s.Shutdown(ctx)
	}()

	slog.Info("Starting server", slog.String("addr", s.Addr))
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
