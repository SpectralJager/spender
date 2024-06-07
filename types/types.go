package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/SpectralJager/spender/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10

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

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}

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

type CreateTimespendParams struct {
	Duration time.Duration `json:"duration"`
	Note     string        `json:"note"`
	Date     time.Time     `json:"date"`
}

func (params CreateTimespendParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.Date.IsZero() {
		errors["date"] = "date should be zero"
	}
	if params.Duration < time.Second {
		errors["duration"] = "duration should be more then 1 second"
	}
	return errors
}

type UpdateTimespendParams struct {
	Duration time.Duration `bson:"duration,omitempty" json:"duration"`
	Date     time.Time     `bson:"date,omitempty" json:"date"`
	Note     string        `bson:"note,omitempty" json:"note"`
}

func (params UpdateTimespendParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.Date.IsZero() {
		errors["date"] = "date should be not zero"
	}
	if params.Duration < time.Second {
		errors["duration"] = "duration should be more then 1 second"
	}
	return errors
}

func (params UpdateTimespendParams) ToBsonDoc() (*bson.D, error) {
	return utils.ToBsonDoc(params)
}

type Timespend struct {
	ID       string        `bson:"_id,omitempty" json:"id,omitempty"`
	OwnerID  string        `bson:"ownerid,omitempty" json:"ownerid,omitempty"`
	Duration time.Duration `bson:"duration" json:"duration"`
	Date     time.Time     `bson:"date" json:"date"`
	Note     string        `bson:"note" json:"note"`
}

func NewTimespendFromParams(params CreateTimespendParams) Timespend {
	return Timespend{
		Duration: params.Duration,
		Date:     params.Date,
		Note:     params.Note,
	}
}

type CreateMoneyspendParams struct {
	Money float64   `json:"money"`
	Note  string    `json:"note"`
	Date  time.Time `json:"date"`
}

func (params CreateMoneyspendParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.Date.IsZero() {
		errors["date"] = "date should be not zero"
	}
	if params.Money <= 0 {
		errors["money"] = "money should be more then 0"
	}
	return errors
}

type UpdateMoneyspendParams struct {
	Money float64   `bson:"money,omitempty" json:"money"`
	Date  time.Time `bson:"date,omitempty" json:"date"`
	Note  string    `bson:"note,omitempty" json:"note"`
}

func (params UpdateMoneyspendParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.Date.IsZero() {
		errors["date"] = "date should be zero"
	}
	if params.Money < 0 {
		errors["money"] = "money should be more positive number"
	}
	return errors
}

func (params UpdateMoneyspendParams) ToBsonDoc() (*bson.D, error) {
	return utils.ToBsonDoc(params)
}

type Moneyspend struct {
	ID      string    `bson:"_id,omitempty" json:"id,omitempty"`
	OwnerID string    `bson:"ownerid,omitempty" json:"ownerid,omitempty"`
	Money   float64   `bson:"money" json:"money"`
	Date    time.Time `bson:"date" json:"date"`
	Note    string    `bson:"note" json:"note"`
}

func NewMoneyspendFromParams(params CreateMoneyspendParams) Moneyspend {
	return Moneyspend{
		Money: params.Money,
		Date:  params.Date,
		Note:  params.Note,
	}
}
