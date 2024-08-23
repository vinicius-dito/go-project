package user_service

import (
	"context"
	"errors"
	"go-project/domain"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type UserService struct {
	userRepository Repository
	tracer         trace.Tracer
}

type UserServiceInput struct {
	UserRepository Repository
	Tracer         trace.Tracer
}

type Repository interface {
	Save(ctx context.Context, user domain.User) error
	Get(ctx context.Context, userID string) (domain.User, error)
}

func NewUserService(input UserServiceInput) (UserService, error) {
	if input.UserRepository == nil {
		return UserService{}, errors.New("missing UserRepository dependency")
	}

	if input.Tracer == nil {
		return UserService{}, errors.New("missing Tracer dependency")
	}

	return UserService{
		userRepository: input.UserRepository,
		tracer:         input.Tracer,
	}, nil
}

func (us UserService) Get(ctx context.Context, input GetDTO) (domain.User, error) {
	ctx, span := us.tracer.Start(ctx, "user-service.Get")
	defer span.End()

	var userDB domain.User

	userDB, err := us.userRepository.Get(ctx, input.UserId)

	return userDB, err
}

func (us UserService) Save(ctx context.Context, input SaveDTO) (domain.User, error) {
	ctx, span := us.tracer.Start(ctx, "user-service.Save")
	defer span.End()

	userDB := domain.User{
		UserId:    input.UserId,
		UserName:  input.UserName,
		Address:   input.Address,
		Birthday:  input.Birthday,
		CreatedAt: time.Now().Format("2006-01-02T15:04:05-0700"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05-0700"),
	}

	// Exemplo de maneira para saber o trace do repository
	// sem precisar inicializo dentro do pkg
	ctx, spanRepository := us.tracer.Start(ctx, "user-repository.Save")
	err := us.userRepository.Save(ctx, userDB)
	spanRepository.End()

	return userDB, err
}
