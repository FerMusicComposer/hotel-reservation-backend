package models

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 8
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

func isEmailValid(email string) bool {
	emailregex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailregex.MatchString(email)
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}

	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("first name must be at least %d characters long", minFirstNameLen))
	}
	if len(params.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("last name must be at least %d characters long", minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password must be at least %d characters long", minPasswordLen))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, "invalid email address")
	}
	return errors
}
