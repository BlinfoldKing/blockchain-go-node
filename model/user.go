package model

import (
	"encoding/json"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/satori/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User use to save user detail
type User struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;"`
	Name         string          `json:"name"`
	NIK          string          `json:"nik"`
	Role         proto.User_Role `json:"role"`
	Username     string          `json:"username"`
	PasswordHash string          `json:"password_hash"`
}

// UserFromJSON use to create user from json string
func UserFromJSON(str string) (user User, err error) {
	err = json.Unmarshal([]byte(str), &user)

	return
}

func (u User) toJSON() (string, error) {
	data, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (u User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func NewUser(
	id uuid.UUID,
	name string,
	nik string,
	role proto.User_Role,
	username string,
	password string,
) (User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}
	return User{
		id,
		name,
		nik,
		role,
		username,
		string(bytes),
	}, nil
}
