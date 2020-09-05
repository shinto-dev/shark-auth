package user

import (
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"shark-auth/pkg/apperrors"
)

func Create(userRepo Repository, userName string, password string) {
	user := User{
		UserId:    uuid.NewV4().String(),
		UserName:  userName,
		Password:  hashPassword(password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)
}

func ExistsByUserName(userRepo Repository, userName string) (bool, error) {
	user, err := userRepo.Get(userName)
	if err != nil {
		return false, errors.Wrap(err, "error while checking user exists")
	}

	return user != (User{}), nil
}

func GetByUserNameAndPassword(userRepo Repository, userName string, password string) (User, error) {
	user, err := userRepo.Get(userName)
	if err != nil {
		return User{}, errors.Wrap(err, "failed while retrieving user")
	}

	if user == (User{}) {
		return User{}, apperrors.ErrUserNotFound
	}

	if !passwordMatch(user.Password, password) {
		return User{}, apperrors.ErrPasswordMismatch
	}
	return user, err
}

func passwordMatch(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	// todo handle this error
	return string(hashedPassword)
}
