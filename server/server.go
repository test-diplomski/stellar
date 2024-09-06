package server

import (
	"fmt"
	sPb "github.com/c12s/scheme/stellar"
	s "github.com/c12s/stellar-go"
	"github.com/c12s/stellar/model"
	"github.com/c12s/stellar/storage"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Server struct {
	db         storage.DB
	instrument map[string]string
}

func (s *Server) List(ctx context.Context, req *sPb.ListReq) (*sPb.ListResp, error) {
	return s.db.List(ctx, req)
}

func (s *Server) Get(ctx context.Context, req *sPb.GetReq) (*sPb.GetResp, error) {
	return s.db.Get(ctx, req)
}

func Run(conf *model.Config, db storage.DB) {
	lis, err := net.Listen("tcp", conf.Address)
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	stellarServer := &Server{
		db:         db,
		instrument: conf.InstrumentConf,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db.StartCollector(ctx)

	n, err := s.NewCollector(stellarServer.instrument["address"], stellarServer.instrument["stopic"])
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := s.InitCollector(stellarServer.instrument["location"], n)
	if err != nil {
		fmt.Println(err)
		return
	}
	go c.Start(ctx, 15*time.Second)

	fmt.Println("StellarService RPC Started")
	sPb.RegisterStellarServiceServer(server, stellarServer)
	server.Serve(lis)
}
