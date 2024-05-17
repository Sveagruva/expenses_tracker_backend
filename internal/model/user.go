package model

type UserModel struct {
	Id           int64
	Login        string
	PasswordHash string
}
