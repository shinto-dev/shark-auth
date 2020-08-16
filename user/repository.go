package user

import (
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	IsValid(userName string, password string) bool
}

type SampleUserRepository map[string]string

func NewSampleRepository(users map[string]string) SampleUserRepository {
	// todo clone it
	return users
}

func (repo SampleUserRepository) IsValid(userName string, password string) bool {
	expectedPassword, ok := repo[userName]

	if !ok || expectedPassword != password {
		logrus.Error("password does not match")
		return false
	}

	return true
}
