package repository

import (
	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/satori/uuid"
)

func (repo *databaseRepository) SaveBalance(balance model.Balance) (err error) {
	err = repo.Save(&balance).Error

	return
}

func (repo *databaseRepository) GetBalanceByUserID(id uuid.UUID) (balance model.Balance, err error) {
	err = repo.DB.Where("user_id = ?", id).Find(&balance).Error

	return
}
