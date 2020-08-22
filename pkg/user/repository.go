package user

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	Get(userName string, password string) (User, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepositoryImpl{
		db: db,
	}
}

func (u UserRepositoryImpl) Get(userName string, password string) (User, error) {
	// todo add encryption for passwords
	row := u.db.QueryRowx("select user_id, user_name from users where user_name=$1 AND password=$2", userName, password)
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
