package main

import (
	"net"
	"os"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-node/server"
	"github.com/blinfoldking/blockchain-go-node/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logrus.SetReportCaller(true)
	service.ServiceConnection = service.New()
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

	handler := server.InitHandler()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.POST("/graphql", handler.Query)
	e.GET("/graphql", handler.Playground)
	e.Logger.Fatal(e.Start(":3000"))

	logrus.Fatal(<-grpcErr)
}
