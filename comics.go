package comics_service

import "errors"

type Comic struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	UserId      int    `json:"user_id" db:"user_id"`
	Month       string `json:"month" db:"month"`
	Day         string `json:"day" db:"day"`
	Year        string `json:"year" db:"year"`
	Img         string `json:"img" db:"img"`
	Description string `json:"description" db:"description"`
}

type UpdateComicInput struct {
	Title       *string `json:"title"`
	Img         *string `json:"img"`
	Description *string `json:"description"`
}

func (i UpdateComicInput) Validate() error {
	if i.Img == nil && i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
