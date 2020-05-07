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

	var serveErr chan error
	go func() {
		blockchainServer := server.Init()
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

		serveErr <- err
	}()

	logrus.Fatal(<-serveErr)
}
