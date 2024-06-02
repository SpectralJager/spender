package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 10
	minFirstNameLen = 2
	maxFirstNameLen = 24
	minLastNameLen  = 2
	maxLastNameLen  = 24
	minPasswordLen  = 8
	maxPasswordLen  = 24
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p CreateUserParams) Validate() []string {
	var errors []string
	if len(p.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("first name lenght should be at least %d characters", minFirstNameLen))
	}
	if len(p.FirstName) > maxFirstNameLen {
		errors = append(errors, fmt.Sprintf("first name lenght should be less or equal then %d characters", maxFirstNameLen))
	}
	if len(p.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("last name lenght should be at least %d characters", maxLastNameLen))
	}
	if len(p.LastName) > maxLastNameLen {
		errors = append(errors, fmt.Sprintf("last name lenght should be less or equal then %d characters", maxLastNameLen))
	}
	if len(p.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password lenght should be at least %d characters", minPasswordLen))
	}
	if len(p.Password) > maxPasswordLen {
		errors = append(errors, fmt.Sprintf("password lenght should be less or equal then %d characters", maxPasswordLen))
	}
	if !emailRegex.MatchString(p.Email) {
		errors = append(errors, "incorect email format")
	}
	return errors
}

type User struct {
	ID           string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName    string `bson:"firstName" json:"firstName"`
	LastName     string `bson:"lastName" json:"lastName"`
	Email        string `bson:"email" json:"email"`
	HashPassword string `bson:"hpassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		HashPassword: string(encpw),
	}, nil
}
