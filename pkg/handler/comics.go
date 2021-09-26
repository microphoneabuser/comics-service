package handler

import (
	"html/template"
	"net/http"
)

var (
	tFeed = template.Must(template.ParseFiles("views/feed.html"))
)

type feedData struct {
	Comics []comicData
}

type comicData struct {
	Id          int
	Title       string
	User        string
	Date        string
	Img         string
	Description string
	IsMine      bool
	IsEdit      bool
}

func (h *Handler) feedGetHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		comics, _ := h.repos.Comics.GetAll()
		data := feedData{}
		for _, comic := range comics {
			user, _ := h.repos.Comics.GetAuthorById(comic.Id)
			data.Comics = append(data.Comics, comicData{
				Id:          comic.Id,
				Title:       comic.Title,
				User:        user.Name,
				Date:        comic.Day + "." + comic.Month + "." + comic.Year,
				Img:         comic.Img,
				Description: comic.Description,
			})
		}
		tFeed.Execute(w, data)
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}

func (h *Handler) userFeedGetHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		comics, _ := h.repos.Comics.GetUsersAll(h.user.Id)
		data := feedData{}
		for _, comic := range comics {
			data.Comics = append(data.Comics, comicData{
				Id:          comic.Id,
				Title:       comic.Title,
				Date:        comic.Day + "." + comic.Month + "." + comic.Year,
				Img:         comic.Img,
				Description: comic.Description,
			})
		}
		tFeed.Execute(w, data)
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}
