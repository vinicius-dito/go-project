package user_service

import (
	"context"
	"errors"
	"go-project/domain"
	"time"
)

type UserService struct {
	userRepository Repository
}

type UserServiceInput struct {
	UserRepository Repository
}

type Repository interface {
	Save(ctx context.Context, user domain.User) error
	Get(ctx context.Context, userID string) (domain.User, error)
}

func NewUserService(input UserServiceInput) (UserService, error) {
	if input.UserRepository == nil {
		return UserService{}, errors.New("missing UserRepository dependency")
	}

	return UserService{
		userRepository: input.UserRepository,
	}, nil
}

func (us UserService) Get(ctx context.Context, input GetDTO) (domain.User, error) {
	var userDB domain.User

	userDB, err := us.userRepository.Get(ctx, input.UserId)

	return userDB, err
}

func (us UserService) Save(ctx context.Context, input SaveDTO) (domain.User, error) {
	userDB := domain.User{
		UserId:    input.UserId,
		UserName:  input.UserName,
		Address:   input.Address,
		Birthday:  input.Birthday,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05-0700"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05-0700"),
	}

	err := us.userRepository.Save(ctx, userDB)

	return userDB, err
}
