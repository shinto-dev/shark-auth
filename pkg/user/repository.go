package user

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	Get(userName string) (User, error)
	Create(user User) error
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepositoryImpl{
		db: db,
	}
}

func (u UserRepositoryImpl) Create(user User) error {
	_, err := u.db.NamedExec("insert into users (user_id, user_name, password, created_at, updated_at)"+
		" values (:user_id, :user_name, :password, :created_at, :updated_at)", user)
	return err
}

func (u UserRepositoryImpl) Get(userName string) (User, error) {
	// todo add unique constraint for userName
	row := u.db.QueryRowx("select user_id, user_name, password from users where user_name=$1", userName)
	if row.Err() != nil {
		logrus.WithError(row.Err()).
			Error("error while querying DB")
		return User{}, row.Err()
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
