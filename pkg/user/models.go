package user

type User struct {
	UserId   string `db:"user_id"`
	UserName string `db:"user_name"`
}
