package signupuser

import (
	"shark-auth/pkg/apperror"
	"shark-auth/pkg/user"
)

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
