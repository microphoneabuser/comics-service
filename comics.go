package comics_service

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
	UserId      int
	Title       string `json:"title"`
	Img         string `json:"img"`
	Description string `json:"description"`
}
