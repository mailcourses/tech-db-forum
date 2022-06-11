package userPostgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
)

type UserRepo struct {
	sqlx *sqlx.DB
}

func NewUserRepo(sqlx *sqlx.DB) domain.UserRepo {
	return UserRepo{sqlx: sqlx}
}

func (repo UserRepo) SelectByID(id int64) (*domain.User, error) {
	query := `SELECT * FROM Users
         	  WHERE id = $1`
	holder := domain.User{}
	if err := repo.sqlx.Get(&holder, query, id); err != nil {
		return nil, err
	}
	return &holder, nil

}

func (repo UserRepo) SelectByNickname(nickname string) (*domain.User, error) {
	query := `SELECT * FROM Users
              WHERE nickname = $1`
	holder := domain.User{}
	if err := repo.sqlx.Get(&holder, query, nickname); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo UserRepo) Create(user domain.User) (*domain.User, error) {
	query := `INSERT INTO Users (nickname, fullname, about, email)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id`
	lastInsertedId := int64(0)
	if err := repo.sqlx.QueryRow(query, user.Nickname, user.Fullname, user.About, user.Email).Scan(&lastInsertedId); err != nil {
		return nil, err
	}

	return repo.SelectByID(lastInsertedId)
}

func (repo UserRepo) Update(user *domain.User) (*domain.User, error) {
	query := `UPDATE Users
			  SET nickname = $1, fullname = $2, about = $3, email = $4
              WHERE id = $5`
	if _, err := repo.sqlx.Exec(query, user.Nickname, user.Fullname, user.About, user.Email, user.Id); err != nil {
		return nil, err
	}

	updatedUser, err := repo.SelectByID(user.Id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
