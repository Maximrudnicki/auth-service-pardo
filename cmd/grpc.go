package main

import (
	"auth_service/cmd/config"
	"auth_service/cmd/model"
	"auth_service/cmd/utils"
	pb "auth_service/proto"
	"context"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GRPCServer) Login(ctx context.Context, users *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Find username in database
	new_users, users_err := s.UsersRepository.FindByEmail(users.Email)
	if users_err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid username or Password",
		)
	}

	config, _ := config.LoadConfig(".")

	verify_error := utils.VerifyPassword(new_users.Password, users.Password)
	if verify_error != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid username or Password",
		)
	}

	// Generate Token
	token, err_token := utils.GenerateToken(config.TokenExpiresIn, new_users.Id, config.TokenSecret)
	if err_token != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot generte token",
		)
	}
	return &pb.LoginResponse{
		TokenType: "Bearer",
		Token:     token}, nil
}

func (s *GRPCServer) Register(ctx context.Context, users *pb.RegisterRequest) (*emptypb.Empty, error) {
	hashedPassword, err := utils.HashPassword(users.Password)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot hash password",
		)
	}

	newUser := model.Users{
		Username: users.Username,
		Email:    users.Email,
		Password: hashedPassword,
	}

	save_err := s.UsersRepository.Save(newUser)
	if save_err != nil {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"Please use another email address",
		)
	}

	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) GetUserId(ctx context.Context, in *pb.TokenRequest) (*pb.UserIdResponse, error) {
	loadConfig, err_conf := config.LoadConfig(".")
	if err_conf != nil {
		log.Fatal("🚀 Could not load token secret", err_conf)
	}

	user, err := utils.ValidateToken(in.Token, loadConfig.TokenSecret)
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"Invalid token",
		)
	}

	userId, err_id := strconv.Atoi(fmt.Sprint(user))
	log.Printf("userId: %v\n", userId)

	if err_id != nil {
		log.Panicf("Failed to listen: %v\n", err_id)
	}

	return &pb.UserIdResponse{
		UserId: uint32(userId),
	}, nil
}
