package repository

import (
	"github.com/jmoiron/sqlx"
	serv "github.com/microphoneabuser/comics_service"
)

type Authorization interface {
	CreateUser(user serv.User) (int, error)
	GetUser(login, password string) (serv.User, error)
	GetUserById(id int) (serv.User, error)
}

type Comics interface {
	Create(comic serv.Comic) (int, error)
	GetAll() ([]serv.Comic, error)
	GetUsersAll(userId int) ([]serv.Comic, error)
	GetAuthorById(comicId int) (serv.User, error)
	GetById(comicId int) (serv.Comic, error)
	Delete(comicId int) error
	Update(comicId int, input serv.UpdateComicInput) error
}

type Repository struct {
	Authorization
	Comics
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Comics:        NewComicsPostgres(db),
	}
}
