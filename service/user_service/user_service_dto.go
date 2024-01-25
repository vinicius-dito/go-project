package user_service

type GetDTO struct {
	UserId string `json:"user_id"`
}

type SaveDTO struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Address  string `json:"address"`
	Birthday string `json:"birthday"`
}
