package user

type User struct {
	ID       string `db:"user_id"`
	Phone    string `db:"phone"`
	Password string `db:"Password"`
}
