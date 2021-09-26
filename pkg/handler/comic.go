package handler

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/microphoneabuser/comics_service"
)

var (
	tComic = template.Must(template.ParseFiles("views/comic.html"))
)

func (h *Handler) comicGetHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		data := comicData{}
		keys := r.URL.Query()
		comicId, err := strconv.Atoi(keys.Get("id"))
		h.comicId = comicId
		if err != nil {
			fmt.Fprintf(w, `404 page not found`)
			return
		}
		comic, err := h.repos.Comics.GetById(comicId)
		if err != nil {
			fmt.Fprintf(w, `404 page not found`)
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
			edit := keys.Get("edit")
			if edit != "" {
				data.IsEdit, err = strconv.ParseBool(edit)
				if err != nil {
					fmt.Fprintf(w, `404 page not found`)
					return
				}
			}
		} else {
			user, _ := h.repos.Comics.GetAuthorById(comicId)
			data.User = user.Name
		}
		tComic.Execute(w, data)
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}

func (h *Handler) comicPostHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		data := comics_service.UpdateComicInput{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
		}
		if data.Title != "" && data.Description != "" {
			var filename string

			keys := r.URL.Query()
			comicId, _ := strconv.Atoi(keys.Get("id"))
			r.ParseMultipartForm(10 << 20)
			file, _, err := r.FormFile("myFile")
			if err == nil {
				defer file.Close()
				if filename, err = h.loadFile(w, r, file); err != nil {
					h.comicGetHandler(w, r)
					fmt.Fprint(w, filename)
					log.Println(err)
					return
				}
				data.Img = filename
			} else {
				comic, err := h.repos.Comics.GetById(comicId)
				if err != nil {
					fmt.Fprintf(w, `404 page not found`)
					return
				}
				data.Img = comic.Img
			}
			err = h.repos.Comics.Update(comicId, data)
			if err != nil {
				h.comicGetHandler(w, r)
				fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
					Неверный формат введеных данных
					</div></div>`)
				log.Printf("%s. comicId: %d\n", err, comicId)
			} else {
				http.Redirect(w, r, fmt.Sprintf("/comic?id=%d", comicId), http.StatusFound)
			}
		} else {
			h.comicGetHandler(w, r)
			fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
				Остались пустые поля!
				</div></div>`)
		}
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}

func (h *Handler) loadFile(w http.ResponseWriter, r *http.Request, file multipart.File) (string, error) {
	tempFile, err := ioutil.TempFile("./views/images", "upload-*.png")
	if err != nil {
		return `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
		Ошибка при загрузке файла.
		</div></div>`, err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
		Ошибка при загрузке файла.
		</div></div>`, err
	}
	filename := strings.TrimPrefix(tempFile.Name(), "./views/")
	tempFile.Write(fileBytes)

	return filename, nil
}

func (h *Handler) comicDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		keys := r.URL.Query()
		comicId, _ := strconv.Atoi(keys.Get("id"))
		if h.comicId != comicId {
			http.Redirect(w, r, "/feed", http.StatusFound)
			return
		}
		if user, _ := h.repos.Comics.GetAuthorById(comicId); user.Id != h.user.Id {
			fmt.Fprintf(w, `404 page not found`)
			return
		}
		if err := h.repos.Comics.Delete(comicId); err != nil {
			fmt.Fprintf(w, `404 page not found`)
			return
		} else {
			http.Redirect(w, r, "/my", http.StatusFound)
		}
	}
}
