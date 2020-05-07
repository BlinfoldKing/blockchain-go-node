package repository

import (
	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/satori/uuid"
)

func (repo databaseRepository) SaveUser(user model.User) (err error) {
	err = repo.DB.Save(&user).Error

	return
}

func (repo databaseRepository) GetUserByID(id uuid.UUID) (user model.User, err error) {
	err = repo.DB.Where("id = ?", id).Find(&user).Error

	return
}
