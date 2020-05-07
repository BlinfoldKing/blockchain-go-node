package server

import (
	"context"
	"time"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/repository"
	"github.com/satori/uuid"
	"github.com/sirupsen/logrus"
)

// Server is implementation of GRPC service
type Server struct {
	repo repository.Repository
}

// Init use to init
func Init() proto.BlockchainServiceServer {
	return &Server{
		repository.Init(),
	}
}

// Ping use to test connection
func (s Server) Ping(ctx context.Context, empty *proto.Empty) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		Ok: true,
	}, nil
}

// CreateUser use to create a new user block
func (s Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.Block, error) {
	userid, err := uuid.FromString(req.Data.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	user := model.User{
		ID:   userid,
		Name: req.Data.GetName(),
		NIK:  req.Data.GetNik(),
	}

	id, err := uuid.FromString(req.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	newBlock, err := model.GenerateNewBlock(
		id,
		req.GetTimestamp(),
		proto.Block_CREATE_USER,
		req.GetPrevHash(),
		user,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &proto.Block{
		Id:        newBlock.ID.String(),
		Timestamp: newBlock.Timestamp.Format(time.RFC3339),
		Nonce:     newBlock.Nonce,
		PrevHash:  newBlock.PrevHash,
		Data:      newBlock.Data,
		Hash:      newBlock.Hash,
	}, nil
}

// MutateBalance create a new balance block
func (s Server) MutateBalance(ctx context.Context, req *proto.MutateBalanceRequest) (*proto.Block, error) {
	bid, err := uuid.FromString(req.Mutation.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	userid, err := uuid.FromString(req.Mutation.GetUserId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	balance := model.Balance{
		ID:       bid,
		UserID:   userid,
		Mutation: req.Mutation.GetMutation(),
	}

	id, err := uuid.FromString(req.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	newBlock, err := model.GenerateNewBlock(
		id,
		req.GetTimestamp(),
		proto.Block_CREATE_USER,
		req.GetPrevHash(),
		balance,
	)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &proto.Block{
		Id:        newBlock.ID.String(),
		Timestamp: newBlock.Timestamp.Format(time.RFC3339),
		Nonce:     newBlock.Nonce,
		BlockType: newBlock.BlockType,
		PrevHash:  newBlock.PrevHash,
		Data:      newBlock.Data,
		Hash:      newBlock.Hash,
	}, nil
}

// PublishBlock use to publish block to this node storage
func (s Server) PublishBlock(ctx context.Context, block *proto.Block) (*proto.Block, error) {
	id, err := uuid.FromString(block.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	timestamp, err := time.Parse(time.RFC3339, block.GetTimestamp())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	newBlock := model.Block{
		ID:        id,
		Timestamp: timestamp,
		Nonce:     block.GetNonce(),
		BlockType: block.GetBlockType(),
		PrevHash:  block.GetPrevHash(),
		Data:      block.GetData(),
		Hash:      block.GetHash(),
	}

	if newBlock.BlockType == proto.Block_CREATE_USER {
		user, err := model.UserFromJSON(newBlock.Data)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		err = s.repo.SaveUser(user)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	} else if newBlock.BlockType == proto.Block_MUTATE_BALANCE {
		balance, err := model.BalanceFromJSON(newBlock.Data)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		err = s.repo.SaveBalance(balance)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	err = s.repo.SaveBlock(newBlock)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return block, nil
}

// QueryBlockchain use to query all blocks
func (s Server) QueryBlockchain(ctx context.Context, req *proto.QueryBlockchainRequest) (*proto.Blockchain, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	blocks, err := s.repo.GetAllBlock(offset, limit)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	blockchain := make([]*proto.Block, 0)
	for _, block := range blocks {
		blockchain = append(blockchain, &proto.Block{
			Id:        block.ID.String(),
			Timestamp: block.Timestamp.Format(time.RFC3339),
			Nonce:     block.Nonce,
			BlockType: block.BlockType,
			PrevHash:  block.PrevHash,
			Data:      block.Data,
			Hash:      block.Hash,
		})
	}

	return &proto.Blockchain{
		Blockchain: blockchain,
	}, nil
}
