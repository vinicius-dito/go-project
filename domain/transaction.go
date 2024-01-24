package domain

type Transaction struct {
	UserId        string
	TransactionId int
	StoreId       int
	SellerId      string
	Revenue       float64
	CreatedAt     string
	UpdatedAt     string
}
