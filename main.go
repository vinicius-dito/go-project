package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	user_firestore_repository "go-project/repository"
	"go-project/service/user_service"
	"io"
	"net/http"
	"os"

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
	bigQueryDataset          string
	bigQueryUsersTable       string
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

	firestoreUserCollection := firestoreClient.Collection(envVars.firestoreUsersCollection)

	userRepository := user_firestore_repository.NewUsersFirestoreRepository(firestoreUserCollection)

	userService, err := user_service.NewUserService(user_service.UserServiceInput{
		UserRepository: userRepository,
	})
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	err = setupHandlerHttp(router, userRepository, userService, ctx)
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
		bigQueryDataset:          os.Getenv("BIG_QUERY_DATASET"),
		bigQueryUsersTable:       os.Getenv("BIG_QUERY_USERS_TABLE"),
	}
}

func setupHandlerHttp(router *mux.Router, userRepository user_firestore_repository.UserFirestoreRepositoy, userService user_service.UserService, ctx context.Context) error {
	router.HandleFunc("/ping", ping)
	//router.HandleFunc("/transaction", transaction)
	//router.HandleFunc("/seller", seller)
	router.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		user(w, req, userRepository, userService, ctx)
	})

	return nil
}

/* func getUser(userRepository users_firestore_repository.UsersFirestoreRepositoy, envVars config) {
	return func(w http.ResponseWriter, req *http.Request) {}
} */

func ping(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "pong")
}

func user(w http.ResponseWriter, req *http.Request, userRepository user_firestore_repository.UserFirestoreRepositoy, userService user_service.UserService, ctx context.Context) {
	switch req.Method {
	case http.MethodGet:
		var user user_service.GetDTO

		getBody, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, fmt.Errorf("%s: %v", read_body_failed, err).Error(), http.StatusBadRequest)
		}

		err = json.Unmarshal(getBody, &user)
		if err != nil {
			http.Error(w, fmt.Errorf("%s: %v", unmarshal_failed, err).Error(), http.StatusUnprocessableEntity)
			return
		}

		//qual código de erro usar aqui?
		if user.UserId == "" {
			http.Error(w, errors.New("empty user_id on body request").Error(), http.StatusInternalServerError)
			return
		}

		userDB, err := userService.Get(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User fetched successfully\n \tuser_id: %v\n \tuser_name: %v\n \taddress: %v\n \tbirthday: %v\n \tcreated_at: %v\n \tupdated_at: %v\n \t",
			userDB.UserId, userDB.UserName,
			userDB.Address, userDB.Birthday,
			userDB.CreatedAt, userDB.UpdatedAt,
		)

	case http.MethodPut:
		var user user_service.SaveDTO

		putBody, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, fmt.Errorf("%s: %v", read_body_failed, err).Error(), http.StatusBadRequest)
		}

		err = json.Unmarshal(putBody, &user)
		if err != nil {
			http.Error(w, fmt.Errorf("%s: %v", unmarshal_failed, err).Error(), http.StatusUnprocessableEntity)
			return
		}

		userDB, err := userService.Save(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User registered successfully\n \t %v", userDB)

	default:
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
	}
}
