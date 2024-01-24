package user_service

import (
	"context"
	"go-project/domain"
	user_firestore_repository "go-project/repository"
)

type UserService struct {
	userRepository user_firestore_repository.UserFirestoreRepositoy
}

type UserServiceInput struct {
	UserRepository user_firestore_repository.UserFirestoreRepositoy
}

func NewUserService(input UserServiceInput) (UserService, error) {
	/* if input.userRepository == nil {
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
	return domain.User{}, nil
}
