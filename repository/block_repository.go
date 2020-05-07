package repository

import "github.com/blinfoldking/blockchain-go-node/model"

func (repo *databaseRepository) SaveBlock(block model.Block) (err error) {
	err = repo.DB.Save(&block).Error

	return
}

func (repo *databaseRepository) GetAllBlock(offset, limit int32) (blocks []model.Block, err error) {
	err = repo.DB.Offset(offset).Limit(limit).Find(&blocks).Error

	return
}
