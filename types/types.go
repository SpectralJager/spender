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

type UpdateUserParams struct {
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
}

func (p UpdateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(p.FirstName) != 0 {
		if len(p.FirstName) < minFirstNameLen {
			errors["firstName"] = fmt.Sprintf("first name lenght should be at least %d characters", minFirstNameLen)
		}
		if len(p.FirstName) > maxFirstNameLen {
			errors["firstName"] = fmt.Sprintf("first name lenght should be less or equal then %d characters", maxFirstNameLen)
		}
	}
	if len(p.LastName) != 0 {
		if len(p.LastName) < minLastNameLen {
			errors["lastName"] = fmt.Sprintf("last name lenght should be at least %d characters", maxLastNameLen)
		}
		if len(p.LastName) > maxLastNameLen {
			errors["lastName"] = fmt.Sprintf("last name lenght should be less or equal then %d characters", maxLastNameLen)
		}
	}
	return errors
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(p.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first name lenght should be at least %d characters", minFirstNameLen)
	}
	if len(p.FirstName) > maxFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first name lenght should be less or equal then %d characters", maxFirstNameLen)
	}
	if len(p.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("last name lenght should be at least %d characters", maxLastNameLen)
	}
	if len(p.LastName) > maxLastNameLen {
		errors["lastName"] = fmt.Sprintf("last name lenght should be less or equal then %d characters", maxLastNameLen)
	}
	if len(p.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password lenght should be at least %d characters", minPasswordLen)
	}
	if len(p.Password) > maxPasswordLen {
		errors["password"] = fmt.Sprintf("password lenght should be less or equal then %d characters", maxPasswordLen)
	}
	if !emailRegex.MatchString(p.Email) {
		errors["email"] = "incorect email format"
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

func NewUserFromParams(params CreateUserParams) (User, error) {
	encpw, err := EncryptPassword(params.Password)
	if err != nil {
		return User{}, err
	}
	return User{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		HashPassword: string(encpw),
	}, nil
}

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}
