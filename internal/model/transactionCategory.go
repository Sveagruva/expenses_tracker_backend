package model

type TransactionCategory struct {
	Id     int64  `json:"id"`
	UserId int64  `json:"userId"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}
