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

	//router := mux.NewRouter()
	router := http.NewServeMux()

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

func setupHandlerHttp(router *http.ServeMux, usersController v1.UserController, ctx context.Context) error {
	router.HandleFunc("GET /ping", func(w http.ResponseWriter, req *http.Request) {

		fmt.Fprint(w, "pong")
	})

	router.HandleFunc("GET /user", func(w http.ResponseWriter, req *http.Request) {

		usersController.GetUser(w, req, ctx)
	})

	/*
		router.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
			usersController.SaveUser(w, req, ctx)
		}).Methods("PUT")
	*/

	return nil
}
