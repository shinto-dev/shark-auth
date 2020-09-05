package user

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	Get(userName string) (User, error)
	Create(user User)
}

type RepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) Repository {
	return RepositoryImpl{
		db: db,
	}
}

func (u RepositoryImpl) Create(user User) {
	_, err := u.db.NamedExec("insert into users (user_id, user_name, password, created_at, updated_at)"+
		" values (:user_id, :user_name, :password, :created_at, :updated_at)", user)
	if err != nil {
		panic(errors.Wrap(err, "error while inserting into DB"))
	}
}

func (u RepositoryImpl) Get(userName string) (User, error) {
	// todo add unique constraint for userName
	row := u.db.QueryRowx("select user_id, user_name, password from users where user_name=$1", userName)
	if row.Err() != nil {
		panic(errors.Wrap(row.Err(), "error while querying DB"))
	}

	var user User
	err := row.StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}

	return user, nil
}
