package service

import (
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/repository"
)

var ServiceConnection Service

type Service struct {
	Repo           repository.Repository
	PoolConnection proto.BlockchainServiceClient
}

func New() Service {
	return Service{
		Repo:           repository.Init(),
		PoolConnection: nil,
	}
}
