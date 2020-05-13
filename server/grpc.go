package server

import (
	"context"
	"time"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/rpc"
	"github.com/blinfoldking/blockchain-go-node/service"
	"github.com/satori/uuid"
	"github.com/sirupsen/logrus"
)

// Server is implementation of GRPC service
type Server struct {
}

// Init use to init
func InitGRPC() proto.BlockchainServiceServer {
	return &Server{}
}

// Ping use to test connection
func (s Server) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	conn, err := rpc.ConnectNode(req.Address)
	if err != nil {
		return nil, err
	}

	service.ServiceConnection.PoolConnection = conn
	return &proto.ConnectResponse{
		Ok: true,
	}, nil
}

// Ping use to test connection
func (s Server) Ping(ctx context.Context, empty *proto.Empty) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		Ok: true,
	}, nil
}

// Count use to count total block
func (s Server) Count(ctx context.Context, empty *proto.Empty) (*proto.BlockCount, error) {
	count, err := service.ServiceConnection.Repo.Count()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &proto.BlockCount{
		Count: count,
	}, nil
}

// GetAllBlock use to count total block
func (s Server) GetAllBlock(ctx context.Context, empty *proto.Empty) (*proto.Blockchain, error) {
	blocks, err := service.ServiceConnection.Repo.GetAllBlock()
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

	count, err := service.ServiceConnection.Repo.Count()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &proto.Blockchain{
		Count:      count,
		Blockchain: blockchain,
	}, nil
}

// CreateUser use to create a new user block
func (s Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.Block, error) {
	logrus.Info("creating block " + req.Data.GetId())
	userid, err := uuid.FromString(req.Data.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	user := model.User{
		ID:           userid,
		Name:         req.Data.GetName(),
		NIK:          req.Data.GetNik(),
		Role:         req.Data.GetRole(),
		Username:     req.Data.GetUsername(),
		PasswordHash: req.Data.GetPasswordHash(),
	}

	id, err := uuid.FromString(req.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	prevBlock, err := service.ServiceConnection.Repo.GetLastBlock()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	newBlock, err := model.GenerateNewBlock(
		id,
		req.GetTimestamp(),
		proto.Block_CREATE_USER,
		prevBlock.Hash,
		user,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	logrus.Info("block created" + userid.String())
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
func (s Server) MutateBalance(ctx context.Context, req *proto.RequestTransaction) (*proto.Block, error) {
	bid, err := uuid.FromString(req.Transaction.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	userid, err := uuid.FromString(req.Transaction.GetUserId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	transaction := model.Transaction{
		ID:     bid,
		UserID: userid,
		Amount: req.Transaction.GetAmount(),
	}

	id, err := uuid.FromString(req.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	prevBlock, err := service.ServiceConnection.Repo.GetLastBlock()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	newBlock, err := model.GenerateNewBlock(
		id,
		req.GetTimestamp(),
		proto.Block_CREATE_USER,
		prevBlock.Hash,
		transaction,
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
		err = service.ServiceConnection.Repo.SaveUser(user)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	} else if newBlock.BlockType == proto.Block_MUTATE_BALANCE {
		balance, err := model.TransactionFromJSON(newBlock.Data)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		err = service.ServiceConnection.Repo.SaveTransaction(balance)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	err = service.ServiceConnection.Repo.SaveBlock(newBlock)
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

	blocks, err := service.ServiceConnection.Repo.QueryAllBlock(offset, limit)
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

	count, err := service.ServiceConnection.Repo.Count()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &proto.Blockchain{
		Count:      count,
		Blockchain: blockchain,
	}, nil
}

func (s Server) GetBlockById(ctx context.Context, req *proto.GetBlockByIdRequest) (*proto.Block, error) {
	id, err := uuid.FromString(req.GetId())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	block, err := service.ServiceConnection.Repo.GetBlockByID(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &proto.Block{
		Id:        block.ID.String(),
		Timestamp: block.Timestamp.Format(time.RFC3339),
		Nonce:     block.Nonce,
		BlockType: block.BlockType,
		PrevHash:  block.PrevHash,
		Data:      block.Data,
		Hash:      block.Hash,
	}, nil
}
