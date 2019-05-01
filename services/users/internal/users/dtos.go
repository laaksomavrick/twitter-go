package users

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (u *User) prepareForInsert() error {
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CompareHashAndPassword(password string) error {
	p := []byte(password)
	hp := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(hp, p)
}
