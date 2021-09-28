package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/microphoneabuser/comics_service"
	"github.com/microphoneabuser/comics_service/utils"
)

var (
	tSignUp = template.Must(template.ParseFiles("views/signup.html"))
)

type signUpInput struct {
	Name     string
	Login    string
	Failed   bool
	ErrorMsg string
}

func (h *Handler) signUpGetHandler(w http.ResponseWriter, r *http.Request) {
	tSignUp.Execute(w, nil)
}

func (h *Handler) signUpPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pass := r.PostForm.Get("password")
	passRep := r.PostForm.Get("password_rep")
	data := signUpInput{}
	if pass != passRep {
		data.Failed = true
		data.ErrorMsg = "Вы ввели разные пароли"
		tSignUp.Execute(w, data)
		return
	}
	hash := utils.GenHash(r.PostForm.Get("password"))
	data.Name = r.PostForm.Get("name")
	data.Login = r.PostForm.Get("login")
	if data.Name != "" && data.Login != "" && pass != "" {
		user := comics_service.User{
			Name:     data.Login,
			Login:    data.Name,
			Password: hash,
		}
		if userId, err := h.repos.Authorization.CreateUser(user); err != nil {
			data.Failed = true
			data.ErrorMsg = "Что-то пошло не так"
			log.Println(err)
			tSignUp.Execute(w, data)
		} else {
			h.user.Id = userId
			http.Redirect(w, r, "/feed", http.StatusFound)
		}
	} else {
		data.Failed = true
		data.ErrorMsg = "Остались пустые поля..."
		tSignUp.Execute(w, data)
	}
}
