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

	return LoginResponse{
		ok:    true,
		token: "",
	}, nil
}
