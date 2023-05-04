package gapi

import (
	"fmt"

	db "github.com/jasonLuFa/simplebank/db/sqlc"
	"github.com/jasonLuFa/simplebank/pb"
	"github.com/jasonLuFa/simplebank/token"
	"github.com/jasonLuFa/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	return server, nil
}
