package data

import (
	"errors"
	"strings"
)

func New() Models {
	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"-"`
}

func (u *User) GetAll() ([]*User, error) {
	usersCopy := make([]*User, len(mockUsers))
	copy(usersCopy, mockUsers)
	return usersCopy, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	for _, user := range mockUsers {
		if strings.EqualFold(user.Email, email) {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (u *User) PasswordMatches(plainText string) (bool, error) {

	if u.Password != plainText {
		return false, errors.New("incorrect password")
	}

	// err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
	// 		// invalid password
	// 		return false, nil
	// 	default:
	// 		return false, err
	// 	}
	// }

	return true, nil
}
