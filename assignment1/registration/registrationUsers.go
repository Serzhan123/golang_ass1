package registration

import (
	"github.com/Bektemis/golang_ass_1/pck"
	"golang.org/x/crypto/bcrypt"
)

type Registration struct {
	Name       string
	SecondName string
	Age        int
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkUsers(name string, d *pck.DatabaseUsers) bool {
	for _, user := range d.Users {
		if user.Name == name {
			return true
		}
	}
	return false
}

func Register(name string, password string, d *pck.DatabaseUsers) bool {
	pass, err := HashPassword(password)
	if err != nil || checkUsers(name, d) {
		return false
	}
	d.Users = append(d.Users, pck.User{Name: name, Password: pass})
	return true
}
