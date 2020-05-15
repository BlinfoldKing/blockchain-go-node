package repository

import (
	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/satori/uuid"
)

func (repo *databaseRepository) SaveTransaction(transaction model.Transaction) (err error) {
	err = repo.Save(&transaction).Error

	return
}

func (repo *databaseRepository) GetTransactionByUserID(id uuid.UUID) (transaction []model.Transaction, err error) {
	err = repo.DB.Where("user_id = ?", id).Find(&transaction).Error

	return
}

func (repo *databaseRepository) GetTransactionByID(id uuid.UUID) (transaction model.Transaction, err error) {
	err = repo.DB.Where("user_id = ?", id).Find(&transaction).Error

	return
}
