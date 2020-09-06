package signupuser

import (
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/user"
)

func CreateUser(userRepo user.Repository, userDetail User) error {
	exists, err := user.ExistsByUserName(userRepo, userDetail.UserName)
	if err != nil {
		return err
	}

	if exists {
		return apperrors.ErrUserNameNotAvailable
	}

	return user.Create(userRepo, userDetail.UserName, userDetail.Password)
}