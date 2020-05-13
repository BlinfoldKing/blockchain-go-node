package model

import (
	"encoding/json"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/satori/uuid"
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
