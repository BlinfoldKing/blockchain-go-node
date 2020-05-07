package repository

import (
	"fmt"
	"os"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/jinzhu/gorm"
	"github.com/satori/uuid"
)

type Repository interface {
	SaveBlock(block model.Block) error
	GetAllBlock(offset, limit int32) ([]model.Block, error)

	SaveUser(user model.User) error
	GetUserByID(id uuid.UUID) (model.User, error)

	SaveBalance(balance model.Balance) error
	GetBalanceByUserID(id uuid.UUID) (model.Balance, error)
}

type databaseRepository struct {
	*gorm.DB
}

func Init() Repository {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.Block{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Balance{})

	return &databaseRepository{db}
}
