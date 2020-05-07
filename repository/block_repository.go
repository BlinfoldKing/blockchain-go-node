package repository

import (
	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/jinzhu/gorm"
)

func (repo *databaseRepository) SaveBlock(block model.Block) (err error) {
	err = repo.DB.Save(&block).Error

	return
}

func (repo *databaseRepository) GetAllBlock() (blocks []model.Block, err error) {
	err = repo.DB.Find(&blocks).Error

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
	err = repo.DB.Count(&count).Error

	return
}
