package main

import (
	"fmt"
	"net"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/server"
	"github.com/joho/godotenv"
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

	fmt.Println(server.Serve(listen))
}
