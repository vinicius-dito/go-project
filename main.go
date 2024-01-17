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

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
)

const (
	method_not_allowed = "METHOD NOT ALLOWED"
	read_body_failed   = "failed to read request body"
	unmarshal_failed   = "failed to unmarshal request body"
)

type config struct {
	gcpProjectId             string
	gcpProjectLocation       string
	firestoreUsersCollection string
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
	UserId    string `firestore:"user_id"`
	UserName  string `firestore:"user_name"`
	Address   string `firestore:"address"`
	Birthday  string `firestore:"birthday"`
	CreatedAt string `firestore:"created_at"`
	UpdatedAt string `firestore:"updated_at"`
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

	// Google Cloud

	firebase, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: envVars.gcpProjectId})
	if err != nil {
		panic(err)
	}

	firestoreClient, err := firebase.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	defer firestoreClient.Close()

	/* bigQueryClient, err := bigquery.NewClient(ctx, envVars.gcpProjectId)
	if err != nil {
		panic(err)
	}
	defer bigQueryClient.Close() */

	router := mux.NewRouter()

	err = handler(router, firestoreClient, envVars, ctx)
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
		gcpProjectId:             os.Getenv("GCP_PROJECT_ID"),
		gcpProjectLocation:       os.Getenv("GCP_PROJECT_LOCATION"),
		firestoreUsersCollection: os.Getenv("FIRESTORE_USERS_COLLECTION"),
	}
}

func handler(router *mux.Router, firestoreClient *firestore.Client, envVars config, ctx context.Context) error {
	// assim eu só aceito método post pra rota. daí eu não conseguia printar o erro.
	//router.HandleFunc("/ping", ping).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/ping", ping)
	router.HandleFunc("/transaction", transaction)
	router.HandleFunc("/seller", seller)
	router.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		user(w, req, firestoreClient, envVars, ctx)
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

func user(w http.ResponseWriter, req *http.Request, firestoreClient *firestore.Client, envVars config, ctx context.Context) {
	switch req.Method {
	case http.MethodGet:
		user, err := getUser(req, firestoreClient, envVars, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User fetched successfully\n \tuser_id: %v\n \tuser_name: %v\n \taddress: %v\n \tbirthday: %v\n \tcreated_at: %v\n \tupdated_at: %v\n \t",
			user.UserId, user.UserName,
			user.Address, user.Birthday,
			user.CreatedAt, user.UpdatedAt,
		)
	case http.MethodPut:
		user, err := saveUser(req, firestoreClient, envVars, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User registered successfully\n \t %v", user)
	default:
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
	}
}

func saveUser(req *http.Request, firestoreClient *firestore.Client, envVars config, ctx context.Context) (users, error) {
	var user usersInput
	var userDB users

	putBody, err := io.ReadAll(req.Body)
	if err != nil {
		return users{}, fmt.Errorf("%s: %v", read_body_failed, err)
	}

	err = json.Unmarshal(putBody, &user)
	if err != nil {
		return users{}, fmt.Errorf("%s: %v", unmarshal_failed, err)
	}

	userDB.Address = user.Address
	userDB.Birthday = user.Birthday
	userDB.UpdatedAt = time.Now().Format("2006-01-02T15:04:05-0700")
	userDB.UserId = user.UserId
	userDB.UserName = user.UserName

	_, err = firestoreClient.Collection(envVars.firestoreUsersCollection).Doc(userDB.UserId).Set(ctx, userDB)

	if err != nil {
		return users{}, fmt.Errorf("failed to insert user into Firestore: %v", err)
	}

	return userDB, nil
}

func getUser(req *http.Request, firestoreClient *firestore.Client, envVars config, ctx context.Context) (users, error) {
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

	userDoc, err := firestoreClient.Collection(envVars.firestoreUsersCollection).Doc(user.UserId).Get(ctx)
	if err != nil {
		return users{}, fmt.Errorf("failed to get user from Firestore: %v", err)
	}

	var userDB users

	if err = userDoc.DataTo(&userDB); err != nil {
		return users{}, fmt.Errorf("failed to parse Firestore document: %v", err)
	}

	return userDB, nil
}
