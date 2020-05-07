package server

import (
	"context"
	"time"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/repository"
	"github.com/satori/uuid"
)

type Server struct {
	repo repository.Repository
}

func Init() proto.BlockchainServiceServer {
	return &Server{
		repository.Init(),
	}
}

func (s Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.Block, error) {
	userid, _ := uuid.FromString(req.Data.GetId())
	user := model.User{
		ID:   userid,
		Name: req.Data.GetName(),
		NIK:  req.Data.GetNik(),
	}

	id, _ := uuid.FromString(req.GetId())
	newBlock, err := model.GenerateNewBlock(
		id,
		req.GetTimestamp(),
		proto.Block_CREATE_USER,
		req.GetPrevHash(),
		user,
	)

	if err != nil {
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

func (s Server) MutateBalance(ctx context.Context, req *proto.MutateBalanceRequest) (*proto.Block, error) {
	bid, _ := uuid.FromString(req.Mutation.GetId())
	userid, _ := uuid.FromString(req.Mutation.GetUserId())
	balance := model.Balance{
		ID:       bid,
		UserID:   userid,
		Mutation: req.Mutation.GetMutation(),
	}

	id, _ := uuid.FromString(req.GetId())
	newBlock, err := model.GenerateNewBlock(
		id,
		req.GetTimestamp(),
		proto.Block_CREATE_USER,
		req.GetPrevHash(),
		balance,
	)

	if err != nil {
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

func (s Server) PublishBlock(ctx context.Context, block *proto.Block) (*proto.Block, error) {
	id, _ := uuid.FromString(block.GetId())
	timestamp, _ := time.Parse(time.RFC3339, block.GetTimestamp())
	newBlock := model.Block{
		ID:        id,
		Timestamp: timestamp,
		Nonce:     block.GetNonce(),
		BlockType: block.GetBlockType(),
		PrevHash:  block.GetPrevHash(),
		Data:      block.GetData(),
		Hash:      block.GetHash(),
	}

	err := s.repo.SaveBlock(newBlock)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (s Server) QueryBlockchain(ctx context.Context, req *proto.QueryBlockchainRequest) (*proto.Blockchain, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	blocks, err := s.repo.GetAllBlock(offset, limit)
	if err != nil {
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
