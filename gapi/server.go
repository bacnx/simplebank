package gapi

import (
	"fmt"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/pb"
	"github.com/bacnx/simplebank/token"
	"github.com/bacnx/simplebank/util"
	"github.com/bacnx/simplebank/worker"
)

// Server serves gRPC requests for our backing service.
type Server struct {
	pb.UnimplementedSimplebankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
