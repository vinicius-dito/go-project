package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

const (
	method_not_allowed = "METHOD NOT ALLOWED"
	read_body_failed   = "failed to read request body"
	unmarshal_failed   = "failed to unmarshal request body"
)

type config struct {
	gcpProjectId       string
	gcpProjectLocation string
	bigQueryDataset    string
	bigQueryUsersTable string
}

type transactions struct {
	userId        string
	transactionId int
	storeId       int
	sellerId      string
	revenue       float64
	createdAt     string
	updatedAt     string
}

type usersInput struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Address  string `json:"address"`
	Birthday string `json:"birthday"`
}

type users struct {
	UserId    string `bigquery:"user_id"`
	UserName  string `bigquery:"user_name"`
	Address   string `bigquery:"address"`
	Birthday  string `bigquery:"birthday"`
	CreatedAt string `bigquery:"created_at"`
	UpdatedAt string `bigquery:"updated_at"`
}

type sellers struct {
	seller_id  string
	store_id   string
	name       string
	document   string
	status     string
	created_at string
	updated_at string
}

func main() {
	ctx := context.Background()

	envVars := setupEnv()

	bigQueryClient, err := bigquery.NewClient(ctx, envVars.gcpProjectId)
	if err != nil {
		panic(err)
	}
	defer bigQueryClient.Close()

	router := mux.NewRouter()

	err = handler(router, bigQueryClient, envVars, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server running in http://localhost:8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Failed to initialize server:", err)
	}
}

func setupEnv() config {
	return config{
		gcpProjectId:       os.Getenv("GCP_PROJECT_ID"),
		gcpProjectLocation: os.Getenv("GCP_PROJECT_LOCATION"),
		bigQueryDataset:    os.Getenv("BIG_QUERY_DATASET"),
		bigQueryUsersTable: os.Getenv("BIG_QUERY_USERS_TABLE"),
	}
}

func handler(router *mux.Router, bigQueryClient *bigquery.Client, envVars config, ctx context.Context) error {
	// assim eu só aceito método post pra rota. daí eu não conseguia printar o erro.
	//router.HandleFunc("/ping", ping).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/ping", ping)
	router.HandleFunc("/transaction", transaction)
	router.HandleFunc("/seller", seller)
	router.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		user(w, req, bigQueryClient, envVars, ctx)
	})

	return nil
}

func ping(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "pong")
}

func transaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// fazer o get
	case http.MethodPost:
		// fazer o post
	default:
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
	}
}

func seller(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// fazer o get
	case http.MethodPost:
		// fazer o post
	default:
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
	}
}

func user(w http.ResponseWriter, req *http.Request, bigQueryClient *bigquery.Client, envVars config, ctx context.Context) {
	switch req.Method {
	case http.MethodGet:
		user, err := getUser(req, bigQueryClient, envVars, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User fetched successfully\n \tuser_id: %v\n \tuser_name: %v\n \taddress: %v\n \tbirthday: %v\n \tcreated_at: %v\n \tupdated_at: %v\n \t",
			user.UserId, user.UserName,
			user.Address, user.Birthday,
			user.CreatedAt, user.UpdatedAt,
		)
	case http.MethodPost:
		user, err := registerUser(req, bigQueryClient, envVars, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User registered successfully\n \t %v", user)
	default:
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
	}
}

func registerUser(req *http.Request, bigQueryClient *bigquery.Client, envVars config, ctx context.Context) (users, error) {
	// fazer uma lógica de dar get no usuário antes de tentar cria-lo.
	var user usersInput
	var userDB users

	postBody, err := io.ReadAll(req.Body)
	if err != nil {
		return users{}, fmt.Errorf("%s: %v", read_body_failed, err)
	}

	err = json.Unmarshal(postBody, &user)
	if err != nil {
		return users{}, fmt.Errorf("%s: %v", unmarshal_failed, err)
	}

	userDB.Address = user.Address
	userDB.Birthday = user.Birthday
	userDB.CreatedAt = time.Now().Format("2006-01-02T15:04:05-0700")
	userDB.UpdatedAt = time.Now().Format("2006-01-02T15:04:05-0700")
	userDB.UserId = user.UserId
	userDB.UserName = user.UserName

	inserter := bigQueryClient.Dataset(envVars.bigQueryDataset).Table(envVars.bigQueryUsersTable).Inserter()

	if err = inserter.Put(ctx, userDB); err != nil {
		return users{}, fmt.Errorf("failed to insert user into BigQuery: %v", err)
	}

	return userDB, nil
}

func getUser(req *http.Request, bigQueryClient *bigquery.Client, envVars config, ctx context.Context) (users, error) {
	var user usersInput

	getBody, err := io.ReadAll(req.Body)
	if err != nil {
		return users{}, fmt.Errorf("%s: %v", read_body_failed, err)
	}

	err = json.Unmarshal(getBody, &user)
	if err != nil {
		return users{}, fmt.Errorf("%s: %v", unmarshal_failed, err)
	}

	if user.UserId == "" {
		return users{}, errors.New("missing user_id on body request")
	}

	queryGet := bigQueryClient.Query(fmt.Sprintf("SELECT * FROM `%s.%s.%s` WHERE user_id = '%s'", envVars.gcpProjectId, envVars.bigQueryDataset, envVars.bigQueryUsersTable, user.UserId))

	queryJob, err := queryGet.Run(ctx)
	if err != nil {
		return users{}, fmt.Errorf("failed to run get user query: %v", err)
	}

	status, err := queryJob.Wait(ctx)
	if err != nil {
		return users{}, fmt.Errorf("failed to retrieve get user query job status: %v", err)
	}

	if err = status.Err(); err != nil {
		return users{}, fmt.Errorf("failed to finish get user query job successfully: %v", err)
	}

	it, err := queryJob.Read(ctx)
	if err != nil {
		return users{}, fmt.Errorf("failed to read get user query job result: %v", err)
	}

	var result users

	for {
		err := it.Next(&result)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return users{}, err
		}
	}

	return result, nil
}
