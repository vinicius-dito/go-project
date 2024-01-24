package user_bigquery_repository

import (
	"context"
	"go-project/domain"

	"cloud.google.com/go/bigquery"
)

type UserBigQueryRepositoy struct {
	bigQueryClient *bigquery.Client
}

func NewUsersBigQueryRepository(bigQueryClient *bigquery.Client) UserBigQueryRepositoy {
	return UserBigQueryRepositoy{
		bigQueryClient: bigQueryClient,
	}
}

func (ubqr UserBigQueryRepositoy) Save(ctx context.Context, user domain.User) error {
	return nil
}

func (ubqr UserBigQueryRepositoy) Get(ctx context.Context, userID string) (domain.User, error) {
	var user domain.User

	return user, nil
}
