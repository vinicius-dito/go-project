package main

import (
	"context"
	"fmt"
	"go-project/config"
	v1 "go-project/handlers/http/v1"
	server "go-project/pkg"
	user_firestore_repository "go-project/repository/firestore"
	"go-project/service/user_service"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
)

const (
	method_not_allowed = "METHOD NOT ALLOWED"
)

func main() {
	ctx := context.Background()

	envVars := config.LoadServerConfig(ctx)

	// GCP Components

	firebase, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: envVars.GCPProjectId})
	if err != nil {
		panic(err)
	}

	firestoreClient, err := firebase.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	defer firestoreClient.Close()

	firestoreUserCollection := firestoreClient.Collection(envVars.FirestoreUsersCollection)

	userRepositoryFirestore := user_firestore_repository.NewUsersFirestoreRepository(firestoreUserCollection)

	//client, err := bigquery.NewClient(ctx, "")

	//userRepositoryBQ := user_bigquery_repository.NewUsersBigQueryRepository(client, "", "", "")

	userService, err := user_service.NewUserService(user_service.UserServiceInput{
		UserRepository: userRepositoryFirestore,
		//UserRepository: userRepositoryBQ,
	})
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	usersController := v1.NewUserController(userService)

	err = setupHandlerHttp(router, usersController, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server running in http://localhost:8080")

	err = server.NewServer(server.ServerInput{
		Port:    ":8080",
		Routers: router,
	})
	if err != nil {
		panic(fmt.Errorf("failed to initialize server: %v", err))
	}
}

func setupHandlerHttp(router *mux.Router, usersController v1.UserController, ctx context.Context) error {
	// go atualizou o pkg http e da pra fazer o router com ele agora.
	router.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
			return
		}

		fmt.Fprint(w, "pong")
	})

	router.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		usersController.GetUser(w, req, ctx)
	}).Methods("GET")

	router.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		usersController.SaveUser(w, req, ctx)
	}).Methods("PUT")

	return nil
}
