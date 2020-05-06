package model

import (
	"encoding/json"

	"github.com/satori/uuid"
)

// User use to save user detail
type User struct {
	ID   uuid.UUID `json:"id" gorm:"primary_key"`
	Name string    `json:"name"`
	NIK  string    `json:"nik"`
}

func (u User) toJSON() (string, error) {
	data, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
