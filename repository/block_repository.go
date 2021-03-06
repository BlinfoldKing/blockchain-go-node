package repository

import (
	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/jinzhu/gorm"
	"github.com/satori/uuid"
)

func (repo *databaseRepository) SaveBlock(block model.Block) (err error) {
	err = repo.DB.Save(&block).Error

	return
}

func (repo *databaseRepository) GetAllBlock() (blocks []model.Block, err error) {
	err = repo.DB.Find(&blocks).Error
	return
}

func (repo *databaseRepository) GetBlockByID(id uuid.UUID) (block model.Block, err error) {
	err = repo.DB.Where("id = ?", id).Find(&block).Error

	return
}

func (repo *databaseRepository) MutateBlockByID(block model.Block) (err error) {
	err = repo.DB.Save(&block).Error

	return
}

func (repo *databaseRepository) QueryAllBlock(offset, limit int32) (blocks []model.Block, err error) {
	err = repo.DB.Offset(offset).Limit(limit).Find(&blocks).Error

	return
}

func (repo *databaseRepository) GetLastBlock() (blocks model.Block, err error) {
	err = repo.DB.Last(&blocks).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return model.Block{}, nil
	}

	return
}

func (repo *databaseRepository) Count() (count int32, err error) {
	err = repo.DB.Model(&model.Block{}).Count(&count).Error

	return
}

func (repo *databaseRepository) DangerouslyDroppingEverything() (err error) {
	err = repo.DB.Delete(&model.Block{}).Error
	if err != nil {
		return
	}
	err = repo.DB.Delete(&model.User{}).Error
	if err != nil {
		return
	}
	err = repo.DB.Delete(&model.Transaction{}).Error

	return
}
