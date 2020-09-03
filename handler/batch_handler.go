package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/toky03/toky-finance-accounting-service/grpc_users"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/service"
	"google.golang.org/grpc"
)

type userServiceBatch interface {
	ReadAllUsers() ([]model.ApplicationUserDTO, model.TokyError)
	UpdateUser(model.ApplicationUserDTO) model.TokyError
	CreateUser(model.ApplicationUserDTO) model.TokyError
	DeleteUser(userId string) model.TokyError
}

type UserServiceServerImpl struct {
	cachedUsers *grpc_users.GetUsersResponse
	userService userServiceBatch
}

func CreateAndRegisterUserBatchService() error {
	userBatchPort := os.Getenv("USER_BATCH_PORT")
	if userBatchPort == "" {
		return errors.New("USER_BATCH_PORT must be specified")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", userBatchPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	grpcServer := grpc.NewServer()
	batchService := UserServiceServerImpl{
		cachedUsers: nil,
		userService: service.CreateApplicationUserService(),
	}
	grpc_users.RegisterUserServiceServer(grpcServer, batchService)
	go grpcServer.Serve(lis)
	return nil
}

func (s UserServiceServerImpl) GetAllUsers(context.Context, *grpc_users.Empty) (*grpc_users.GetUsersResponse, error) {
	if s.cachedUsers != nil && len(s.cachedUsers.GetUsers()) > 0 {
		return s.cachedUsers, nil
	}
	users, tokyErr := s.userService.ReadAllUsers()
	if model.IsExisting(tokyErr) {
		return nil, tokyErr.Error()
	}
	s.cachedUsers = mapUserDTOsTogrpcUsers(users)
	return s.cachedUsers, nil
}
func (s UserServiceServerImpl) UpdateUser(ctx context.Context, user *grpc_users.User) (*grpc_users.Empty, error) {
	tokyErr := s.userService.UpdateUser(mapGrpcUserToDTO(user))
	if model.IsExisting(tokyErr) {
		return nil, tokyErr.Error()
	}
	s.cachedUsers = nil
	return &grpc_users.Empty{}, nil
}
func (s UserServiceServerImpl) AddUser(ctx context.Context, user *grpc_users.User) (*grpc_users.Empty, error) {
	tokyErr := s.userService.CreateUser(mapGrpcUserToDTO(user))
	if model.IsExisting(tokyErr) {
		return nil, tokyErr.Error()
	}
	s.cachedUsers = nil
	return &grpc_users.Empty{}, nil
}
func (s UserServiceServerImpl) DeleteUser(ctx context.Context, userId *grpc_users.UserId) (*grpc_users.Empty, error) {
	tokyErr := s.userService.DeleteUser(userId.GetId())
	if model.IsExisting(tokyErr) {
		return nil, tokyErr.Error()
	}
	s.cachedUsers = nil
	return &grpc_users.Empty{}, nil
}

func mapUserDTOsTogrpcUsers(usersDTO []model.ApplicationUserDTO) *grpc_users.GetUsersResponse {
	usersGrpc := make([]*grpc_users.User, 0, len(usersDTO))
	for _, userDTO := range usersDTO {
		usersGrpc = append(usersGrpc, &grpc_users.User{
			Id:        userDTO.UserID,
			Email:     userDTO.EMail,
			Firstname: userDTO.FirstName,
			Lastname:  userDTO.LastName,
			Username:  userDTO.UserName})
	}
	return &grpc_users.GetUsersResponse{Users: usersGrpc}
}

func mapGrpcUserToDTO(grpcUser *grpc_users.User) model.ApplicationUserDTO {
	return model.ApplicationUserDTO{
		UserID:    grpcUser.GetId(),
		EMail:     grpcUser.GetEmail(),
		FirstName: grpcUser.GetFirstname(),
		LastName:  grpcUser.GetLastname(),
		UserName:  grpcUser.GetUsername()}
}
