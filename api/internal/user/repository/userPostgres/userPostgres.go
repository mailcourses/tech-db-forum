package userPostgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
	"golang.org/x/net/context"
	"strings"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) domain.UserRepo {
	return UserRepo{pool: pool}
}

func (repo UserRepo) SelectById(id int64) (*domain.User, error) {
	query := `SELECT nickname, fullname, about, email FROM Users
         	  WHERE id = $1`
	holder := domain.User{}

	if err := repo.pool.QueryRow(context.Background(), query, id).Scan(
		&holder.Nickname,
		&holder.Fullname,
		&holder.About,
		&holder.Email); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo UserRepo) SelectByNickname(nickname string) (*domain.User, error) {
	query := `SELECT nickname, fullname, about, email FROM Users
              WHERE nickname = $1`
	holder := domain.User{}
	if err := repo.pool.QueryRow(context.Background(), query, strings.ToLower(nickname)).Scan(
		&holder.Nickname,
		&holder.Fullname,
		&holder.About,
		&holder.Email); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo UserRepo) Create(user domain.User) ([]domain.User, error) {
	preQuery := `select nickname, fullname, about, email from users where nickname=$1 or email=$2`
	rows, err := repo.pool.Query(context.Background(), preQuery, strings.ToLower(user.Nickname), strings.ToLower(user.Email))
	if err != nil {
		return nil, err
	}

	var checked []domain.User
	for rows.Next() {
		element := domain.User{}
		if err := rows.Scan(
			&element.Nickname,
			&element.Fullname,
			&element.About,
			&element.Email); err != nil {
			return nil, err
		}
		checked = append(checked, element)
	}
	if len(checked) > 0 {
		return checked, &userErrors.UserErrorConfilct{Conflict: "nickname or email"}
	}

	query := `
			INSERT INTO Users (nickname, fullname, about, email)
	  		VALUES ($1, $2, $3, $4)
   		    returning nickname, fullname, about, email;`

	inserted := domain.User{}

	if err := repo.pool.QueryRow(context.Background(), query, user.Nickname, user.Fullname, user.About, user.Email).Scan(domain.GetUserFields(&inserted)...); err != nil {
		return nil, err
	}

	return []domain.User{inserted}, nil
}

func (repo UserRepo) Update(user *domain.User) (*domain.User, error) {
	checked := domain.User{}
	preQuery := `select nickname, fullname, about, email from users where lower(email)=$1`
	if err := repo.pool.QueryRow(context.Background(), preQuery, strings.ToLower(user.Email)).Scan(&checked); err == nil &&
		strings.ToLower(checked.Nickname) != strings.ToLower(user.Nickname) {
		return nil, &userErrors.UserErrorConfilct{Conflict: "email"}
	}

	query := `UPDATE Users
			  SET fullname = $2, about = $3, email = $4
              WHERE nickname = $1
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
	if err := repo.pool.QueryRow(context.Background(), query, strings.ToLower(user.Nickname), user.Fullname, user.About, user.Email).Scan(
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
		query = `SELECT DISTINCT nickname, fullname, about, email FROM ForumUsers u
			 WHERE u.forum = $1 AND u.nickname > $2
			 ORDER BY nickname DESC
			 LIMIT $3`
	} else {
		query = `SELECT DISTINCT nickname, fullname, about, email FROM ForumUsers u
			 WHERE u.forum = $1 AND u.nickname > $2
			 ORDER BY nickname
			 LIMIT $3`
	}

	rows, err := repo.pool.Query(context.Background(), query, slug, strings.ToLower(since), limit)
	if err != nil {
		return nil, err
	}

	var users []domain.User
	for rows.Next() {
		element := domain.User{}
		if err := rows.Scan(
			&element.Nickname,
			&element.Fullname,
			&element.About,
			&element.Email); err != nil {
			return nil, err
		}
		users = append(users, element)
	}

	return users, nil

}
