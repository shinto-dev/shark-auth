package user

import (
	"time"
)

type User struct {
	UserId    string    `db:"user_id"`
	UserName  string    `db:"user_name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
