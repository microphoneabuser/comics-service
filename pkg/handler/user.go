package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/microphoneabuser/comics_service"
	"github.com/microphoneabuser/comics_service/utils"
)

var (
	tUser = template.Must(template.ParseFiles("views/user.html"))
)

type userInfo struct {
	Id       int
	Name     string
	Login    string
	Password string
	Comics   []comicData
	IsMine   bool
	IsEdit   bool
	PassEdit bool
}

func (h *Handler) userGetHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/user?id=%d", h.user.Id), http.StatusFound)
		}
		user, _ := h.repos.Authorization.GetUserById(id)
		data := userInfo{
			Id:       id,
			Name:     user.Name,
			Login:    user.Login,
			Password: user.Password,
		}

		if h.user.Id == id {
			data.IsMine = true
		}
		edit, _ := strconv.Atoi(keys.Get("edit"))
		if edit == 1 {
			data.IsEdit = true
		}
		passEdit, _ := strconv.Atoi(keys.Get("pass"))
		if passEdit == 1 {
			data.PassEdit = true
		}

		comics, _ := h.repos.Comics.GetUsersAll(id)
		for _, comic := range comics {
			data.Comics = append(data.Comics, comicData{
				Id:    comic.Id,
				Title: comic.Title,
				Date:  comic.Day + "." + comic.Month + "." + comic.Year,
				Img:   comic.Img,
			})
		}
		tUser.Execute(w, data)
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}

func (h *Handler) userPostHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		if err != nil {
			fmt.Fprintf(w, `404 page not found`)
			return
		}
		edit, err := strconv.Atoi(keys.Get("edit"))
		if edit != 1 {
			fmt.Fprintf(w, `404 page not found`)
			return
		}
		if err != nil {
			fmt.Fprintf(w, `404 page not found`)
			return
		}

		data := comics_service.UpdateUserInput{
			Name:  r.FormValue("name"),
			Login: r.FormValue("login"),
		}

		if data.Name != "" && data.Login != "" {
			newPass, _ := strconv.Atoi(keys.Get("pass"))
			if newPass == 1 {
				pass := r.FormValue("password")
				passRep := r.FormValue("password_rep")
				if pass == "" || passRep == "" {
					fmt.Fprintf(w, `<script>alert("Вы не ввели пароль...")</script>`)
					h.userGetHandler(w, r)
				}
				if pass != passRep {
					fmt.Fprintf(w, `<script>alert("Вы ввели разные пароли!!!")</script>`)
					h.userGetHandler(w, r)
				}
				data.Password = utils.GenHash(pass)
			} else {
				usr, _ := h.repos.Authorization.GetUserById(id)
				data.Password = usr.Password
			}
			if err := h.repos.Authorization.Update(id, data); err != nil {
				fmt.Fprintf(w, `<script>alert("Произошла ошибка при изменении пользователя...")</script>`)
				h.userGetHandler(w, r)
			} else {
				http.Redirect(w, r, fmt.Sprintf("/user?id=%d", h.user.Id), http.StatusFound)
			}
		} else {
			fmt.Fprintf(w, `<script>alert("Остались пустые поля...")</script>`)
			h.userGetHandler(w, r)
		}
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}
