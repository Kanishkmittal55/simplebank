package gapi

import (
	"fmt"
	"github.com/kanishkmittal55/simplebank/worker"

	db "github.com/kanishkmittal55/simplebank/db/sqlc"
	"github.com/kanishkmittal55/simplebank/db/util"
	"github.com/kanishkmittal55/simplebank/pb"
	"github.com/kanishkmittal55/simplebank/token"
)

// Server serves HTTP request for our banking service.
type Server struct {
	pb.UnimplementedSimplebankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new grpc server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token : %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
