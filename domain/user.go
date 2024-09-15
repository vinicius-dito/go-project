package domain

type User struct {
	UserId    string `firestore:"user_id"`
	UserName  string `firestore:"user_name"`
	Address   string `firestore:"address"`
	Birthday  string `firestore:"birthday"`
	CreatedAt string `firestore:"created_at"`
	UpdatedAt string `firestore:"updated_at"`
}
