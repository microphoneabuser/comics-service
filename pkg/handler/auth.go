package handler

import (
	"html/template"
	"net/http"

	"github.com/microphoneabuser/comics_service"
	"github.com/microphoneabuser/comics_service/utils"
)

var (
	tAuth = template.Must(template.ParseFiles("views/auth.html"))
)

type signInInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Failed   bool   `json:"success"`
	Error    error  `json:"error"`
}

func (h *Handler) authGetHandler(w http.ResponseWriter, r *http.Request) {
	tAuth.Execute(w, nil)
}

func (h *Handler) authPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := utils.GenHash(r.PostForm.Get("password"))
	data := signInInput{
		Login:    r.PostForm.Get("login"),
		Password: hash,
		Failed:   false,
		Error:    nil,
	}
	if data.Login != "" && data.Password != "" {
		if user, err := h.repos.Authorization.GetUser(data.Login, data.Password); err != nil {
			data.Failed = true
			data.Error = err
			tAuth.Execute(w, data)
		} else {
			data.Failed = false
			h.user = user
			http.Redirect(w, r, "/feed", http.StatusFound)
		}
	} else {
		h.user = comics_service.User{}
		tAuth.Execute(w, data)
	}
}
