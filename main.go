package main

import (
	"net"
	"os"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	var grpcErr chan error
	go func() {
		blockchainServer := server.InitGRPC()
		server := grpc.NewServer()
		proto.RegisterBlockchainServiceServer(server, blockchainServer)

		port := os.Getenv("PORT")
		port = ":" + port
		listen, err := net.Listen("tcp", port)
		if err != nil {
			panic(err)
		}
		logrus.Info("serve on port 9000")
		err = server.Serve(listen)

		grpcErr <- err
	}()

	logrus.Fatal(<-grpcErr)
}
