package handler

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/microphoneabuser/comics_service"
)

var (
	tUpload = template.Must(template.ParseFiles("views/upload.html"))
)

type uploadInput struct {
	Title       string
	Description string
}

func (h *Handler) uploadComicGetHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		tUpload.Execute(w, r)
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}

func (h *Handler) uploadComicPostHandler(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		data := uploadInput{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
		}
		if data.Title != "" && data.Description != "" {
			var filename string
			var err error
			if filename, err = h.loadFile(w, r); err != nil {
				tUpload.Execute(w, data)
				fmt.Fprint(w, filename)
				log.Println(err)
				return
			}
			comicId, err := h.repos.Comics.Create(comics_service.Comic{
				Title:       data.Title,
				UserId:      h.user.Id,
				Month:       fmt.Sprint(int(time.Now().Month())),
				Day:         fmt.Sprint(time.Now().Day()),
				Year:        fmt.Sprint(time.Now().Year()),
				Img:         filename,
				Description: data.Description,
			})
			if err != nil {
				tUpload.Execute(w, data)
				fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
					Неверный формат введеных данных
					</div></div>`)
				log.Printf("%s. comicId: %d\n", err, comicId)
			} else {
				tUpload.Execute(w, data)
				fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
				Успешно!
				</div></div>`)
				http.Redirect(w, r, "/feed", http.StatusFound)
			}
		} else {
			tUpload.Execute(w, data)
			fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
				Остались пустые поля!
				</div></div>`)
		}
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}

func (h *Handler) loadFile(w http.ResponseWriter, r *http.Request) (string, error) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("myFile")
	if err != nil {
		return `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
		Вы не загрузили файл.
		</div></div>`, err
	}
	defer file.Close()

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
