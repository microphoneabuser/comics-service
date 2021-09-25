package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	serv "github.com/microphoneabuser/comics_service"
)

type ComicsPostgres struct {
	db *sqlx.DB
}

func NewComicsPostgres(db *sqlx.DB) *ComicsPostgres {
	return &ComicsPostgres{db: db}
}

func (r *ComicsPostgres) GetAll() ([]serv.Comic, error) {
	var comics []serv.Comic
	query := fmt.Sprintf("SELECT id, title, user_id, date_part('month', date) as month, date_part('day', date) as day, date_part('year', date) as year, img, description FROM %s", comicsTable)
	err := r.db.Select(&comics, query)
	return comics, err
}

func (r *ComicsPostgres) GetUsersAll(userId int) ([]serv.Comic, error) {
	var comics []serv.Comic
	query := fmt.Sprintf(`SELECT id, title, user_id, 
		date_part('month', date) as month, date_part('day', date) as day, 
		date_part('year', date) as year, img, description FROM %s
		WHERE user_id = $1`, comicsTable)
	err := r.db.Select(&comics, query, userId)
	return comics, err
}

func (r *ComicsPostgres) GetAuthorById(comicId int) (serv.User, error) {
	var user serv.User
	query := fmt.Sprintf("SELECT u.id, u.name, u.login FROM %s u INNER JOIN %s c on u.id = c.user_id WHERE c.id=$1", usersTable, comicsTable)
	err := r.db.Get(&user, query, comicId)
	return user, err
}

func (r *ComicsPostgres) Create(comic serv.Comic) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (title, user_id, date, img, description) 
		values ($1, $2, make_date($3, $4, $5), $6, $7) RETURNING id`, comicsTable)
	row := r.db.QueryRow(query, comic.Title, comic.UserId, comic.Year, comic.Month, comic.Day, comic.Img, comic.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ComicsPostgres) GetById(comicId int) (serv.Comic, error) {
	var comic serv.Comic
	query := fmt.Sprintf(`SELECT id, title, user_id, date_part('month', date) as month, date_part('day', date) as day, 
	date_part('year', date) as year, img, description FROM %s WHERE id=$1`, comicsTable)
	err := r.db.Get(&comic, query, comicId)
	return comic, err
}

func (r *ComicsPostgres) Delete(comicId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", comicsTable)
	_, err := r.db.Exec(query, comicId)
	return err
}

func (r *ComicsPostgres) Update(comicId int, input serv.UpdateComicInput) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, img = $2, description = $3 WHERE id = $4", comicsTable)
	_, err := r.db.Exec(query, input.Title, input.Img, input.Description, comicId)
	return err
}
