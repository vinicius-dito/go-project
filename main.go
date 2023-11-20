package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	err := api(router)
	if err != nil {
		panic(err)
	}

	fmt.Println("Servidor rodando em http://localhost:8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func api(router *mux.Router) error {
	// assim eu só aceito método post pra rota. daí eu não conseguia printar o erro.
	//router.HandleFunc("/", handleRoot).Methods(http.MethodPost)
	router.HandleFunc("/", handleRoot)

	return nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Bem-vindo à minha API HTTP simples!")
}
