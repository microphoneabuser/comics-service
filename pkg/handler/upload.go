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
	IsFirst = false
)

type uploadInput struct {
	Title       string
	Description string
}

func (h *Handler) uploadComic(w http.ResponseWriter, r *http.Request) {
	if h.user.Id != 0 {
		tUpload.Execute(w, r)
		data := uploadInput{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
		}
		if !IsFirst {
			IsFirst = true
			return
		}

		if data.Title != "" && data.Description != "" {
			r.ParseMultipartForm(10 << 20)
			file, _, err := r.FormFile("myFile")
			if err != nil {
				fmt.Println("Error Retrieving the File")
				fmt.Println(err)
				fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
				Вы не загрузили файл.
				</div></div>`)
				return
			}
			defer file.Close()

			tempFile, err := ioutil.TempFile("./views/images", "upload-*.png")
			if err != nil {
				log.Println(err.Error() + "(Error Retrieving the File)")
			}
			defer tempFile.Close()

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println(err)
			}
			filename := strings.TrimPrefix(tempFile.Name(), "./views/")
			tempFile.Write(fileBytes)

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
				fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
					Неверный формат введеных данных
					</div></div>`)
				log.Printf("%s. comicId: %d\n", err, comicId)
			} else {
				// http.Redirect(w, r, "/comics?id="+fmt.Sprint(comicId), http.StatusFound)
				fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
				Успешно!
				</div></div>`)
				defer http.Redirect(w, r, "/feed", http.StatusFound)
				IsFirst = false
			}
		} else {
			fmt.Fprintf(w, `<div class="bd justify-content-center"><div class="alert alert-success" role="alert">
				Остались пустые поля!
				</div></div>`)
		}
	} else {
		http.Redirect(w, r, "/auth", http.StatusFound)
	}
}
