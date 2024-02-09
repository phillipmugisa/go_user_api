package data

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/phillipmugisa/go_user_api/types"
)

type User struct {
	UserName      string     `json:"username"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	Phone         string     `json:"phone"`
	Region        string     `json:"region"`
	UserGender    string     `json:"userGender"`
	UserLanguage  string     `json:"userLanguage"`
	UserDateBirth CustomDate `json:"userDateBirth"`
	Code          string     `json:"code"`
}

type CustomDate struct {
	time.Time
}

const customDateFormat = "02.01.2006"

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse(customDateFormat, s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		cd.Time = v
	default:
		return fmt.Errorf("unsupported type %T for CustomDate", value)
	}
	return nil
}

func (u *User) Validate() *types.ApiError {
	// necesssary data checks
	return nil
}

// generates verification code
func (u *User) GenerateCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	part1 := rand.Intn(1000)
	part2 := rand.Intn(1000)

	randomNumber := fmt.Sprintf("%03d-%03d", part1, part2)

	return randomNumber
}
