package model

import (
	"encoding/json"

	"github.com/satori/uuid"
)

// User use to save user detail
type User struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name string    `json:"name"`
	NIK  string    `json:"nik"`
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
