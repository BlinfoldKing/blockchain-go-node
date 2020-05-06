package model

import (
	"encoding/json"

	"github.com/satori/uuid"
)

// Balance use to track money transaction of a user
type Balance struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
	UserID   uuid.UUID `json:"user_id"`
	Mutation int32     `json:"mutation"`
}

func (b Balance) toJSON() (string, error) {
	data, err := json.Marshal(&b)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
