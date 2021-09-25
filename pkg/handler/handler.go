package handler

import (
	"net/http"

	"github.com/microphoneabuser/comics_service"

	"github.com/microphoneabuser/comics_service/pkg/repository"
)

type Handler struct {
	repos *repository.Repository
	user  comics_service.User
}

func NewHandler(repos *repository.Repository) *Handler {
	return &Handler{repos: repos}
}

func (h *Handler) SetupRoutes() {
	http.HandleFunc("/", redirect)
	http.HandleFunc("/auth", h.authHandler)
	http.HandleFunc("/feed", h.feedHandler)
	http.HandleFunc("/my", h.userFeedHandler)
	http.HandleFunc("/upload", h.uploadComic)
	http.HandleFunc("/comic", h.getComic)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("views/images/"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/feed", http.StatusFound)
}
