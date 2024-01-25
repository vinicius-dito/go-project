package user_bigquery_repository

import (
	"context"
	"fmt"
	"go-project/domain"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type UserBigQueryRepositoy struct {
	bigQueryClient    *bigquery.Client
	bigQueryProjectID string
	bigQueryDataset   string
	bigQueryTable     string
}

func NewUsersBigQueryRepository(bigQueryClient *bigquery.Client, bigQueryProjectID string, bigQueryDataset string, bigQueryTable string) UserBigQueryRepositoy {
	return UserBigQueryRepositoy{
		bigQueryClient:    bigQueryClient,
		bigQueryProjectID: bigQueryProjectID,
		bigQueryDataset:   bigQueryDataset,
		bigQueryTable:     bigQueryTable,
	}
}

func (ubqr UserBigQueryRepositoy) Save(ctx context.Context, user domain.User) error {
	inserter := ubqr.bigQueryClient.Dataset(ubqr.bigQueryDataset).Table(ubqr.bigQueryTable).Inserter()

	if err := inserter.Put(ctx, user); err != nil {
		return fmt.Errorf("failed to insert user into BigQuery: %v", err)
	}

	return nil
}

func (ubqr UserBigQueryRepositoy) Get(ctx context.Context, userID string) (domain.User, error) {
	var user domain.User

	queryGet := ubqr.bigQueryClient.Query(fmt.Sprintf("SELECT * FROM `%s.%s.%s` WHERE user_id = '%s'", ubqr.bigQueryProjectID, ubqr.bigQueryDataset, ubqr.bigQueryTable, userID))

	queryJob, err := queryGet.Run(ctx)
	if err != nil {
		return user, fmt.Errorf("failed to run get user query: %v", err)
	}

	status, err := queryJob.Wait(ctx)
	if err != nil {
		return user, fmt.Errorf("failed to retrieve get user query job status: %v", err)
	}

	if err = status.Err(); err != nil {
		return user, fmt.Errorf("failed to finish get user query job successfully: %v", err)
	}

	it, err := queryJob.Read(ctx)
	if err != nil {
		return user, fmt.Errorf("failed to read get user query job result: %v", err)
	}

	for {
		err := it.Next(&user)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return user, err
		}
	}

	return user, nil
}
