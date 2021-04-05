package internal

import (
	"shark-auth/internal/apperror"
	"shark-auth/internal/user"
)

type User struct {
	UserName string
	Password string
}

// todo add country, dob etc.

func CreateUser(userRepo user.Repository, userDetail User) error {
	exists, err := user.ExistsByUserName(userRepo, userDetail.UserName)
	if err != nil {
		return err
	}

	if exists {
		return apperror.NewError(apperror.CodeInvalidRequest, "user name already taken")
	}

	return user.Create(userRepo, userDetail.UserName, userDetail.Password)
}
