package resolver

import (
	"context"
	"errors"
	"time"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/service"
	"github.com/satori/uuid"
	"github.com/sirupsen/logrus"
)

type Resolver struct{}

func New() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Register(ctx context.Context, args struct {
	req struct {
		Name     string
		Nik      string
		Username string
		Password string
	}
}) (res LoginResponse, err error) {
	user, err := model.NewUser(
		uuid.NewV4(),
		args.req.Name,
		args.req.Nik,
		proto.User_CLIENT,
		args.req.Username,
		args.req.Password,
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
			Id:        user.ID.String(),
			Timestamp: time.Now().Format(time.RFC3339),
			Data: &proto.User{
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
