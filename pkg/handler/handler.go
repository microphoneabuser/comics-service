package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/microphoneabuser/comics_service"
	"github.com/microphoneabuser/comics_service/pkg/repository"
)

type Handler struct {
	repos   *repository.Repository
	user    comics_service.User
	comicId int
}

func NewHandler(repos *repository.Repository) *Handler {
	return &Handler{repos: repos}
}

func (h *Handler) SetupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", redirect)
	r.HandleFunc("/auth", h.authGetHandler).Methods("GET")
	r.HandleFunc("/auth", h.authPostHandler).Methods("POST")

	r.HandleFunc("/feed", h.feedGetHandler).Methods("GET")
	r.HandleFunc("/my", h.userFeedGetHandler).Methods("GET")

	r.HandleFunc("/upload", h.uploadComicGetHandler).Methods("GET")
	r.HandleFunc("/upload", h.uploadComicPostHandler).Methods("POST")

	r.HandleFunc("/comic", h.comicGetHandler).Methods("GET")
	r.HandleFunc("/comic", h.comicPostHandler).Methods("POST")

	r.HandleFunc("/del", h.comicDeleteHandler).Methods("GET")

	r.HandleFunc("/user", h.userGetHandler).Methods("GET")
	r.HandleFunc("/user", h.userPostHandler).Methods("POST")

	fsi := http.FileServer(http.Dir("views/images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", fsi))
	fss := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fss))
	http.Handle("/", r)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/feed", http.StatusFound)
}
