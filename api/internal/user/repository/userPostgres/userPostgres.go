package userPostgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
	"strings"
)

type UserRepo struct {
	sqlx *sqlx.DB
}

func NewUserRepo(sqlx *sqlx.DB) domain.UserRepo {
	return UserRepo{sqlx: sqlx}
}

func (repo UserRepo) SelectById(id int64) (*domain.User, error) {
	query := `SELECT nickname, fullname, about, email FROM Users
         	  WHERE id = $1`
	holder := domain.User{}

	if err := repo.sqlx.Get(&holder, query, id); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo UserRepo) SelectByNickname(nickname string) (*domain.User, error) {
	query := `SELECT nickname, fullname, about, email FROM Users
              WHERE lower(nickname) = $1`
	holder := domain.User{}
	if err := repo.sqlx.Get(&holder, query, strings.ToLower(nickname)); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo UserRepo) Create(user domain.User) ([]domain.User, error) {
	var checked []domain.User

	preQuery := `select nickname, fullname, about, email from users where lower(nickname)=$1 or lower(email)=$2`
	if _ = repo.sqlx.Select(&checked, preQuery, strings.ToLower(user.Nickname), strings.ToLower(user.Email)); len(checked) > 0 {
		return checked, &userErrors.UserErrorConfilct{Conflict: "nickname or email"}
	}
	checked = nil
	query := `
			INSERT INTO Users (nickname, fullname, about, email)
	  		VALUES ($1, $2, $3, $4)
   		    returning nickname, fullname, about, email;`

	inserted := domain.User{}

	if err := repo.sqlx.QueryRow(query, user.Nickname, user.Fullname, user.About, user.Email).Scan(
		&inserted.Nickname,
		&inserted.Fullname,
		&inserted.About,
		&inserted.Email); err != nil {
		return nil, err
	}

	return []domain.User{inserted}, nil
}

func (repo UserRepo) Update(user *domain.User) (*domain.User, error) {
	checked := domain.User{}
	preQuery := `select nickname, fullname, about, email from users where lower(email)=$1`
	if err := repo.sqlx.Select(&checked, preQuery, strings.ToLower(user.Email)); err == nil &&
		strings.ToLower(checked.Nickname) != strings.ToLower(user.Nickname) {
		return nil, &userErrors.UserErrorConfilct{Conflict: "email"}
	}

	query := `UPDATE Users
			  SET fullname = $2, about = $3, email = $4
              WHERE lower(nickname) = $1
              RETURNING nickname, fullname, about, email`

	beforeUpd, _ := repo.SelectByNickname(user.Nickname)

	if user.Fullname == "" {
		user.Fullname = beforeUpd.Fullname
	}

	if user.Email == "" {
		user.Email = beforeUpd.Email
	}

	if user.About == "" {
		user.About = beforeUpd.About
	}

	updated := domain.User{}
	if err := repo.sqlx.QueryRow(query, strings.ToLower(user.Nickname), user.Fullname, user.About, user.Email).Scan(
		&updated.Nickname,
		&updated.Fullname,
		&updated.About,
		&updated.Email); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (repo UserRepo) SelectUsersBySlug(slug string, limit int64, since string, desc bool) ([]domain.User, error) {
	var query string

	if desc {
		query = `SELECT DISTINCT nickname, fullname, about, email FROM Users U
			 JOIN Forum f ON f.user_nickname = U.nickname
             JOIN Thread t ON t.user_nickname = U.nickname
             WHERE f.slug = $1 AND lower(nickname) > $2
			 ORDER BY lower(nickname) DESC
			 LIMIT $3`
	} else {
		query = `SELECT DISTINCT nickname, fullname, about, email FROM Users U
			 JOIN Forum f ON f.user_nickname = U.nickname
             JOIN Thread t ON t.user_nickname = U.nickname
             WHERE f.slug = $1 AND lower(nickname) > $2
			 ORDER BY lower(nickname)
			 LIMIT $3`
	}

	var users []domain.User

	if err := repo.sqlx.Select(&users, query, slug, strings.ToLower(since), limit); err != nil {
		return nil, err
	}

	return users, nil

}
