package authorization

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/Bektemis/golang_ass_1/pck"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignIn(user, password string, users *pck.DatabaseUsers) bool {
	for _, u := range users.Users {
		if u.Name == user && CheckPasswordHash(password, u.Password) {
			return true
		}
	}
	return false
}
