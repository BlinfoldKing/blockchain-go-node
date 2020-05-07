package model

import (
	"encoding/json"

	"github.com/satori/uuid"
)

// Balance use to track money transaction of a user
type Balance struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	UserID   uuid.UUID `json:"user_id"`
	Mutation int32     `json:"mutation"`
}

// BalanceFromJSON use to create user from json string
func BalanceFromJSON(str string) (balance Balance, err error) {
	err = json.Unmarshal([]byte(str), &balance)

	return
}

func (b Balance) toJSON() (string, error) {
	data, err := json.Marshal(&b)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
