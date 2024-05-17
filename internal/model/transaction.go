package model

type Transaction struct {
	Id         int64               `json:"id"`
	Price      int64               `json:"price"`
	CategoryId int64               `json:"categoryId"`
	CreatedAt  string              `json:"createdAt"`
	UserId     int64               `json:"userId"`
	Category   TransactionCategory `json:"category"`
}
