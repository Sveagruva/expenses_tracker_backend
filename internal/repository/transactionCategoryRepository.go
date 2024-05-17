package repository

import (
	"database/sql"
	"expenses_tracker/internal/model"
)

type TransactionCategoryRepository interface {
	CreateTransactionCategory(category model.TransactionCategory) error
	GetTransactionCategoryById(categoryId int64) (model.TransactionCategory, error)
	GetTransactionCategories(userId int64) ([]model.TransactionCategory, error)
	DeleteTransactionCategory(id int64) error
	UpdateCategoryById(id int64, name string, color string) error
}

type transactionCategoryRepository struct {
	db *sql.DB
}

func GetTransactionCategoryRepository(db *sql.DB) *transactionCategoryRepository {
	return &transactionCategoryRepository{db: db}
}

func (repo *transactionCategoryRepository) CreateTransactionCategory(category model.TransactionCategory) error {
	query := `INSERT INTO "TransactionCategories" ("UserId", "Name", "Color") VALUES ($1, $2, $3)`
	_, err := repo.db.Exec(query, category.UserId, category.Name, category.Color)
	return err
}

func (repo *transactionCategoryRepository) GetTransactionCategoryById(categoryId int64) (model.TransactionCategory, error) {
	var category model.TransactionCategory
	query := `SELECT "Id", "UserId", "Name", "Color" FROM "TransactionCategories" WHERE "Id" = $1 LIMIT 1`
	err := repo.db.QueryRow(query, categoryId).Scan(&category.Id, &category.UserId, &category.Name, &category.Color)
	return category, err
}

func (repo *transactionCategoryRepository) GetTransactionCategories(userId int64) ([]model.TransactionCategory, error) {
	var categories []model.TransactionCategory
	query := `SELECT "Id", "UserId", "Name", "Color" FROM "TransactionCategories" WHERE "UserId" = $1`
	rows, err := repo.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.TransactionCategory
		if err := rows.Scan(&category.Id, &category.UserId, &category.Name, &category.Color); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (repo *transactionCategoryRepository) DeleteTransactionCategory(id int64) error {
	query := `DELETE FROM "TransactionCategories" WHERE "Id" = $1`
	_, err := repo.db.Exec(query, id)
	return err
}

func (repo *transactionCategoryRepository) UpdateCategoryById(id int64, name string, color string) error {
	query := `UPDATE "TransactionCategories" SET "Name" = $1, "Color" = $2 WHERE "Id" = $3`
	_, err := repo.db.Exec(query, name, color, id)
	return err
}
