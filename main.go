package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-project/domain"
	users_firestore_repository "go-project/repository"
	"io"
	"net/http"
	"os"
	"time"

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

type usersInput struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Address  string `json:"address"`
	Birthday string `json:"birthday"`
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

	usersRepository := users_firestore_repository.NewUsersFirestoreRepository(*firestoreClient, envVars.firestoreUsersCollection)

	router := mux.NewRouter()

	err = handler(router, usersRepository, envVars, ctx)
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

func handler(router *mux.Router, usersRepository users_firestore_repository.UsersFirestoreRepositoy, envVars config, ctx context.Context) error {
	router.HandleFunc("/ping", ping)
	router.HandleFunc("/transaction", transaction)
	router.HandleFunc("/seller", seller)
	router.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		user(w, req, usersRepository, envVars, ctx)
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

func user(w http.ResponseWriter, req *http.Request, usersRepository users_firestore_repository.UsersFirestoreRepositoy, envVars config, ctx context.Context) {
	switch req.Method {
	case http.MethodGet:
		user, err := getUser(req, usersRepository, envVars, ctx)
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
		user, err := saveUser(req, usersRepository, envVars, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User registered successfully\n \t %v", user)
	default:
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
	}
}

func saveUser(req *http.Request, usersRepository users_firestore_repository.UsersFirestoreRepositoy, envVars config, ctx context.Context) (domain.Users, error) {
	var user usersInput
	var userDB domain.Users

	putBody, err := io.ReadAll(req.Body)
	if err != nil {
		return userDB, fmt.Errorf("%s: %v", read_body_failed, err)
	}

	err = json.Unmarshal(putBody, &user)
	if err != nil {
		return userDB, fmt.Errorf("%s: %v", unmarshal_failed, err)
	}

	userDB.Address = user.Address
	userDB.Birthday = user.Birthday
	userDB.UpdatedAt = time.Now().Format("2006-01-02T15:04:05-0700")
	userDB.UserId = user.UserId
	userDB.UserName = user.UserName

	err = usersRepository.Save(ctx, userDB)
	if err != nil {
		return domain.Users{}, err
	}

	return userDB, nil
}

func getUser(req *http.Request, usersRepository users_firestore_repository.UsersFirestoreRepositoy, envVars config, ctx context.Context) (domain.Users, error) {
	var user usersInput
	var userDB domain.Users

	getBody, err := io.ReadAll(req.Body)
	if err != nil {
		return userDB, fmt.Errorf("%s: %v", read_body_failed, err)
	}

	err = json.Unmarshal(getBody, &user)
	if err != nil {
		return userDB, fmt.Errorf("%s: %v", unmarshal_failed, err)
	}

	if user.UserId == "" {
		return userDB, errors.New("missing user_id on body request")
	}

	userDB, err = usersRepository.Get(ctx, user.UserId)
	if err != nil {
		return userDB, err
	}

	return userDB, nil
}
