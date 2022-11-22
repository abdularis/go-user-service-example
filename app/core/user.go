package core

import "golang.org/x/crypto/bcrypt"

type UserRole string

var (
	Admin       UserRole = "admin"
	DefaultUser UserRole = "user"
)

type User struct {
	ID       uint     `json:"id" gorm:"column:id"`
	UserName string   `json:"userName" gorm:"column:username"`
	Password string   `json:"-" gorm:"column:password"`
	Role     UserRole `json:"role" gorm:"column:role"`
}

func (u *User) SetRoleFromStr(role string) {
	switch role {
	case "admin":
		u.Role = Admin
	default:
		u.Role = DefaultUser
	}
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
