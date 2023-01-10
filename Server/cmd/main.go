package main

import (
	"net"

	"Server/internal/connection"
	pb "Server/internal/proto"
	"Server/internal/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {

	log := logrus.New()

	PSQL, err := connection.PSQLConnect()
	if err != nil {
		log.Fatalf("couldn't connect postgres: %s", err)
	}

	Mongo, err := connection.MongoConnect()
	if err != nil {
		log.Fatalf("couldn't connect mongo: %s", err)
	}

	Service := service.NewService(PSQL, Mongo)
	Service.PingPSQL()
	Service.PingMongo()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("FAIL PORT")
	}

	knockServer := pb.NewKnockServer(Service)

	s := grpc.NewServer()
	pb.RegisterKnockingServer(s, knockServer)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("SERVER UPAL")
	}
}
