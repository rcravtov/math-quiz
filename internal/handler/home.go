package handler

import (
	"math-quiz/internal/view/home"
	"net/http"
)

type homeHandler struct {
	BaseURL string
}

func (h homeHandler) handleIndex(w http.ResponseWriter, r *http.Request) error {
	props := home.HomeProps{BaseURL: h.BaseURL}
	return home.Index(props).Render(r.Context(), w)
}
