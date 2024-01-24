package user_service

import (
	"context"
	"go-project/domain"
	user_firestore_repository "go-project/repository"
	"time"
)

type UserService struct {
	userRepository user_firestore_repository.UserFirestoreRepositoy
}

type UserServiceInput struct {
	UserRepository user_firestore_repository.UserFirestoreRepositoy
}

func NewUserService(input UserServiceInput) (UserService, error) {
	// como retornar esse erro aqui?
	/* if input.UserRepository == nil {
		return UserService{}, errors.New("missing userRepository dependency")
	} */

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
