package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var (
	tComic = template.Must(template.ParseFiles("views/comic.html"))
)

func (h *Handler) getComic(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		data := comicData{}
		keys := r.URL.Query()
		comicId, err := strconv.Atoi(keys.Get("id"))
		if err != nil {
			tComic.Execute(w, data)
			fmt.Fprintf(w, `<div class="alert alert-success" role="alert">
				Данной страницы не существует hdfsjgk
				</div>`)
			return
		}
		comic, err := h.repos.Comics.GetById(comicId)
		if err != nil {
			tComic.Execute(w, data)
			fmt.Fprintf(w, `<div class="alert alert-success" role="alert">
				Данной страницы не существует
				</div>`)
			return
		}
		data = comicData{
			Id:          comic.Id,
			Title:       comic.Title,
			Date:        comic.Day + "." + comic.Month + "." + comic.Year,
			Img:         comic.Img,
			Description: comic.Description,
			IsMine:      false,
		}
		if comic.UserId == h.user.Id {
			data.IsMine = true
		} else {
			user, _ := h.repos.Comics.GetAuthorById(comicId)
			data.User = user.Name
		}
		tComic.Execute(w, data)
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}
