package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	IsValid(userName string, password string) (bool, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepositoryImpl{
		db: db,
	}
}

func (u UserRepositoryImpl) IsValid(userName string, password string) (bool, error) {
	// todo add encryption for passwords
	row := u.db.QueryRow("select count(1) from users where user_name=$1 AND password=$2", userName, password)
	if row.Err() != nil {
		logrus.WithError(row.Err()).
			Error("error while querying DB")
		return false, row.Err()
	}

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}
