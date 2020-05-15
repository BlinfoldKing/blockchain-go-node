package resolver

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/rpc"
	"github.com/blinfoldking/blockchain-go-node/service"
	"github.com/satori/uuid"
	"github.com/sirupsen/logrus"
)

type Resolver struct{}

func New() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Connect(ctx context.Context, args struct {
	Url string
}) (bool, error) {
	conn, err := rpc.ConnectNode(args.Url)
	if err != nil {
		return false, err
	}
	service.ServiceConnection.PoolConnection = conn
	conn.Connect(ctx, &proto.ConnectRequest{
		Address: os.Getenv("SELF_URL"),
	})
	return true, nil
}

func (r *Resolver) Register(ctx context.Context, args struct {
	Request struct {
		Name     string
		Nik      string
		Username string
		Password string
	}
}) (res LoginResponse, err error) {
	user, err := model.NewUser(
		uuid.NewV4(),
		args.Request.Name,
		args.Request.Nik,
		proto.User_CLIENT,
		args.Request.Username,
		args.Request.Password,
	)
	if err != nil {
		logrus.Error(err)
		return
	}

	if service.ServiceConnection.PoolConnection == nil {
		err = errors.New("no data pool connected")
		return
	}

	_, err = service.ServiceConnection.
		PoolConnection.
		CreateUser(ctx, &proto.CreateUserRequest{
			Id:        uuid.NewV4().String(),
			Timestamp: time.Now().Format(time.RFC3339),
			Data: &proto.User{
				Id:           user.ID.String(),
				Name:         user.Name,
				Nik:          user.NIK,
				Role:         user.Role,
				Username:     user.Username,
				PasswordHash: user.PasswordHash,
			},
		})

	token, err := user.GenerateToken()
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ok:    true,
		token: token,
	}, nil
}

func (r Resolver) Login(ctx context.Context, args struct {
	Request struct {
		Username string
		Password string
	}
}) (LoginResponse, error) {
	user, err := service.ServiceConnection.Repo.GetUserByUsername(args.Request.Username)
	if err != nil {
		return LoginResponse{ok: false}, err
	}

	ok := user.CheckPasswordHash(args.Request.Password)
	if !ok {
		return LoginResponse{ok: false}, errors.New("wrong password")
	}

	token, err := user.GenerateToken()
	if err != nil {
		return LoginResponse{ok: true}, err
	}

	return LoginResponse{
		ok:    true,
		token: token,
	}, nil
}

func (r Resolver) GetAccountDetail(ctx context.Context) (UserResolver, error) {
	userI := ctx.Value("user")
	if userI == nil {
		return UserResolver{}, errors.New("no valid token found")
	}

	user := userI.(model.User)

	return UserResolver{user}, nil
}

func (r Resolver) MutateBalance(ctx context.Context, args struct {
	Amount int32
}) (TransactionResolver, error) {
	userI := ctx.Value("user")
	if userI == nil {
		return TransactionResolver{}, errors.New("no valid token found")
	}

	user := userI.(model.User)

	req := &proto.RequestTransaction{
		Id:        uuid.NewV4().String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Transaction: &proto.Transaction{
			Id:     uuid.NewV4().String(),
			UserId: user.ID.String(),
			Amount: args.Amount,
		},
	}
	_, err := service.ServiceConnection.PoolConnection.MutateBalance(ctx, req)
	if err != nil {
		return TransactionResolver{}, err
	}
	trId, _ := uuid.FromString(req.Transaction.GetId())

	return TransactionResolver{
		model.Transaction{
			ID:     trId,
			UserID: user.ID,
			Amount: req.Transaction.GetAmount(),
		},
	}, nil
}
