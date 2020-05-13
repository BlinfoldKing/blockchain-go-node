package model

import (
	"encoding/json"

	"github.com/satori/uuid"
)

// Transaction use to track money transaction of a user
type Transaction struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	UserID uuid.UUID `json:"user_id"`
	Amount int32     `json:"mutation"`
}

// TransactionFromJSON use to create user from json string
func TransactionFromJSON(str string) (transaction Transaction, err error) {
	err = json.Unmarshal([]byte(str), &transaction)

	return
}

func (t Transaction) toJSON() (string, error) {
	data, err := json.Marshal(&t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
