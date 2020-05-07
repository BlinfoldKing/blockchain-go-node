package main

import (
	"fmt"
	"net"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	blockchainServer := server.Init()
	server := grpc.NewServer()
	proto.RegisterBlockchainServiceServer(server, blockchainServer)

	var serveErr chan error
	go func() {
		logrus.Info("serve on port 9000")
		err = server.Serve(listen)

		serveErr <- err
	}()

	fmt.Println("hello")
	<-serveErr

	fmt.Println(serveErr)
}
