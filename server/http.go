package server

import "github.com/blinfoldking/blockchain-go-node/server/graphql"

type Handler struct {
	graphql.GraphQLHandler
}

func InitHandler() Handler {
	h := Handler{}

	return h
}
