package repository

import (
	"database/sql"
	"errors"
	"expenses_tracker/internal/model"
)

type UserRepository interface {
	Create(user model.UserModel) error
	FindByLogin(login string) (model.UserModel, error)
	FindById(id int64) (model.UserModel, error)
}

type userRepository struct {
	db *sql.DB
}

func GetUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) Create(user model.UserModel) error {
	query := `INSERT INTO "Users" ("Login", "PasswordHash") VALUES ($1, $2)`
	err := repo.db.QueryRow(query, user.Login, user.PasswordHash).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func (repo *userRepository) FindByLogin(login string) (model.UserModel, error) {
	var user model.UserModel
	query := `SELECT "Id", "Login", "PasswordHash" FROM "Users" WHERE "Login" = $1`
	err := repo.db.QueryRow(query, login).Scan(&user.Id, &user.Login, &user.PasswordHash)
	return user, err
}

func (repo *userRepository) FindById(id int64) (model.UserModel, error) {
	var user model.UserModel
	query := `SELECT "Id", "Login", "PasswordHash" FROM "Users" WHERE "Id" = $1`
	err := repo.db.QueryRow(query, id).Scan(&user.Id, &user.Login, &user.PasswordHash)
	return user, err
}
