package repository

import (
	"database/sql"
	"expenses_tracker/internal/model"
	"expenses_tracker/internal/pkg/utils"
	"fmt"
	"strconv"
	"strings"
)

type TransactionRepository interface {
	CreateTransaction(transaction model.Transaction) error
	GetTransactionById(transactionId int64) (model.Transaction, error)
	GetTransactions(userId int64, categoryIds []int64, pagination SqlPagination) (PaginationResponse[model.Transaction], error)
	UpdateTransaction(transaction model.Transaction) error
	DeleteTransaction(id int64) error
	GetTotalPriceByDateAndCategory(userId int64, year int, month int, day int, categoryId int64) (float64, error)
}

type transactionRepository struct {
	db *sql.DB
}

func GetTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (repo *transactionRepository) CreateTransaction(transaction model.Transaction) error {
	query := `INSERT INTO "Transactions" ("Price", "CategoryId", "UserId") VALUES ($1, $2, $3)`
	_, err := repo.db.Exec(query, transaction.Price, transaction.CategoryId, transaction.UserId)
	return err
}

func (repo *transactionRepository) GetTransactionById(transactionId int64) (model.Transaction, error) {
	var transaction model.Transaction
	query := `SELECT "Id", "Price", "CategoryId", "CreatedAt", "UserId" FROM "Transactions" WHERE "Id" = $1 LIMIT 1`
	err := repo.db.QueryRow(query, transactionId).Scan(&transaction.Id, &transaction.Price, &transaction.CategoryId, &transaction.CreatedAt, &transaction.UserId)
	if err != nil {
		return model.Transaction{}, err
	}
	return transaction, nil
}

func (repo *transactionRepository) GetTransactions(userId int64, categoryIds []int64, pagination SqlPagination) (PaginationResponse[model.Transaction], error) {
	var transactions []model.Transaction = []model.Transaction{}

	counter := &utils.IncreasingCounter{}

	mainQuery := `
        SELECT "Transactions"."Id", "Price", "CategoryId", "CreatedAt", "Transactions"."UserId", cat."Id", cat."name", cat."color"
        FROM "Transactions"
        INNER JOIN "TransactionCategories" as cat on "Transactions"."CategoryId" = cat."Id"
        WHERE "Transactions"."UserId" = ` + "$" + strconv.Itoa(counter.Next())
	queryParams := []interface{}{userId}

	if len(categoryIds) > 0 {
		placeholders := make([]string, len(categoryIds))
		for i, categoryId := range categoryIds {
			placeholders[i] = "$" + strconv.Itoa(counter.Next())
			queryParams = append(queryParams, categoryId)
		}

		mainQuery += " AND \"CategoryId\" IN ("
		mainQuery += strings.Join(placeholders, ", ")
		mainQuery += ")"
	}

	countSubquery := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM (%s) AS subquery`, mainQuery)

	var totalCount int64
	if err := repo.db.QueryRow(countSubquery, queryParams...).Scan(&totalCount); err != nil {
		return PaginationResponse[model.Transaction]{Items: transactions, Count: 0}, err
	}

	mainQuery += " ORDER BY \"Transactions\".\"CreatedAt\" DESC LIMIT $" + strconv.Itoa(counter.Next()) + " OFFSET $" + strconv.Itoa(counter.Next())
	queryParams = append(queryParams, pagination.Limit, pagination.Offset)

	rows, err := repo.db.Query(mainQuery, queryParams...)
	if err != nil {
		return PaginationResponse[model.Transaction]{Items: transactions, Count: 0}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Transaction
		if err := rows.Scan(&item.Id, &item.Price, &item.CategoryId, &item.CreatedAt, &item.UserId, &item.Category.Id, &item.Category.Name, &item.Category.Color); err != nil {
			return PaginationResponse[model.Transaction]{Items: transactions, Count: 0}, err
		}
		transactions = append(transactions, item)
	}

	return PaginationResponse[model.Transaction]{Items: transactions, Count: totalCount}, nil
}

func (repo *transactionRepository) DeleteTransaction(id int64) error {
	query := `DELETE FROM "Transactions" WHERE "Id" = $1`
	_, err := repo.db.Exec(query, id)
	return err
}

func (repo *transactionRepository) UpdateTransaction(transaction model.Transaction) error {
	query := `UPDATE "Transactions" SET "Price" = $1 WHERE "Id" = $2`
	_, err := repo.db.Exec(query, transaction.Price, transaction.Id)
	return err
}

func (repo *transactionRepository) GetTotalPriceByDateAndCategory(userId int64, year int, month int, day int, categoryId int64) (float64, error) {
	var totalPrice float64

	counter := utils.IncreasingCounter{}
	conditions := []string{
		`"UserId" = $` + fmt.Sprintf("%d", counter.Next()),
		`strftime('%Y', "CreatedAt") = $` + fmt.Sprintf("%d", counter.Next()),
	}
	args := []interface{}{userId, fmt.Sprintf("%04d", year)}

	if categoryId != 0 {
		conditions = append(conditions, `"CategoryId" = $`+fmt.Sprintf("%d", counter.Next()))
		args = append(args, categoryId)
	}

	if month != 0 {
		conditions = append(conditions, `strftime('%m', "CreatedAt") = $`+fmt.Sprintf("%d", counter.Next()))
		args = append(args, fmt.Sprintf("%02d", month))
		if day != 0 {
			conditions = append(conditions, `strftime('%d', "CreatedAt") = $`+fmt.Sprintf("%d", counter.Next()))
			args = append(args, fmt.Sprintf("%02d", day))
		}
	}

	query := fmt.Sprintf(`SELECT COALESCE(SUM("Price"), 0) FROM "Transactions" WHERE %s`, strings.Join(conditions, " AND "))

	err := repo.db.QueryRow(query, args...).Scan(&totalPrice)
	if err != nil {
		return 0, err
	}

	return totalPrice, nil
}
