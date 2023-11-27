package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	method_not_allowed = "METHOD NOT ALLOWED"
	read_body_failed   = "failed to read request body"
	unmarshal_failed   = "failed to unmarshal request body"
)

type transactions struct {
	userId        string
	transactionId int
	storeId       int
	sellerId      string
	revenue       float64
	createdAt     string
	updatedAt     string
}

type users struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Address   string `json:"address"`
	Birthday  string `json:"birthday"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type sellers struct {
	seller_id  string
	store_id   int
	name       string
	document   int
	status     string
	created_at string
	updated_at string
}

func main() {
	router := mux.NewRouter()

	err := handler(router)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server running in http://localhost:8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Failed to initialize server:", err)
	}
}

func handler(router *mux.Router) error {
	// assim eu só aceito método post pra rota. daí eu não conseguia printar o erro.
	//router.HandleFunc("/ping", ping).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/ping", ping)
	router.HandleFunc("/transaction", transaction)
	router.HandleFunc("/seller", seller)
	router.HandleFunc("/user", user)

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

func user(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		user, err := getUser(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User fetched successfully\n \tUserId: %v", user.UserId)
	case http.MethodPost:
		user, err := registerUser(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User registered successfully: %v", user)
	default:
		http.Error(w, method_not_allowed, http.StatusMethodNotAllowed)
	}
}

func registerUser(req *http.Request) (users, error) {
	// fazer uma lógica de dar get no usuário antes de tentar cria-lo.
	// fazer separação de domínio, pois não quero que venha created_at e updated_at da api.
	var user users

	postBody, err := io.ReadAll(req.Body)
	if err != nil {
		return users{}, errors.New(read_body_failed)
	}

	err = json.Unmarshal(postBody, &user)
	if err != nil {
		return users{}, errors.New(unmarshal_failed)
	}

	return users{
		UserId:    user.UserId,
		UserName:  user.UserName,
		Address:   user.Address,
		Birthday:  user.Birthday,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}, nil
}

func getUser(req *http.Request) (users, error) {
	var user users

	getBody, err := io.ReadAll(req.Body)
	if err != nil {
		return users{}, errors.New(read_body_failed)
	}

	err = json.Unmarshal(getBody, &user)
	if err != nil {
		return users{}, errors.New(unmarshal_failed)
	}

	if user.UserId == "" {
		return users{}, errors.New("missing user_id on body request")
	}

	return users{
		UserId: user.UserId,
	}, nil
}
