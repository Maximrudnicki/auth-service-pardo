package main

import (
	"auth_service/cmd/config"
	"auth_service/cmd/model"
	pb "auth_service/proto"
	"auth_service/cmd/repository"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.AuthenticationServiceServer
	UsersRepository repository.UsersRepository
}

func main() {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	//Database
	db := config.ConnectionDB(&loadConfig)

	db_table_err := db.Table("users").AutoMigrate(&model.Users{})
	if db_table_err != nil {
		log.Fatalf("Databese table error: %v\n", db_table_err)
	}

	//Init Repository
	userRepository := repository.NewUsersRepositoryImpl(db)

	// Start GRPC Server
	lis, err := net.Listen("tcp", loadConfig.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	log.Printf("Listening at %s\n", loadConfig.GRPCPort)

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	pb.RegisterAuthenticationServiceServer(s, &GRPCServer{UsersRepository: userRepository})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
