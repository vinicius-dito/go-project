package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ApplicationName          string
	TraceRatio               float64
	TraceGCPProjectId        string
	ApplicationEnvironment   string
	GCPProjectId             string
	GCPProjectLocation       string
	FirestoreUsersCollection string
	BigQueryDataset          string
	BigQueryUsersTable       string
}

func LoadServerConfig(ctx context.Context) Config {
	traceRatio, err := strconv.ParseFloat(os.Getenv("TRACE_RATIO"), 64)
	if err != nil {
		panic(fmt.Errorf("failed to parse TRACE_RATIO env var to float: %v", err))
	}

	return Config{
		ApplicationName:          os.Getenv("APPLICATION_NAME"),
		TraceRatio:               traceRatio,
		TraceGCPProjectId:        os.Getenv("TRACE_GCP_PROJECT_ID"),
		ApplicationEnvironment:   os.Getenv("APPLICATION_ENVIRONMENT"),
		GCPProjectId:             os.Getenv("GCP_PROJECT_ID"),
		GCPProjectLocation:       os.Getenv("GCP_PROJECT_LOCATION"),
		FirestoreUsersCollection: os.Getenv("FIRESTORE_USERS_COLLECTION"),
		BigQueryDataset:          os.Getenv("BIG_QUERY_DATASET"),
		BigQueryUsersTable:       os.Getenv("BIG_QUERY_USERS_TABLE"),
	}
}
