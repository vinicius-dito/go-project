package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-project/service/user_service"
	"io"
	"net/http"
)

const (
	read_body_failed = "failed to read request body"
	unmarshal_failed = "failed to unmarshal request body"
)

type UserController struct {
	userService user_service.UserService
}

func NewUserController(userService user_service.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (uc UserController) GetUser(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	var user user_service.GetDTO

	queryParam := req.URL.Query()

	user.UserId = queryParam.Get("user_id")

	if user.UserId == "" {
		http.Error(w, errors.New("empty user_id on query params request").Error(), http.StatusBadRequest)
		return
	}

	userDB, err := uc.userService.Get(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User fetched successfully\n \tuser_id: %v\n \tuser_name: %v\n \taddress: %v\n \tbirthday: %v\n \tcreated_at: %v\n \tupdated_at: %v\n \t",
		userDB.UserId, userDB.UserName,
		userDB.Address, userDB.Birthday,
		userDB.CreatedAt, userDB.UpdatedAt,
	)
}

func (uc UserController) SaveUser(w http.ResponseWriter, req *http.Request, ctx context.Context) {
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

	userDB, err := uc.userService.Save(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User registered successfully\n \t %v", userDB)
}
