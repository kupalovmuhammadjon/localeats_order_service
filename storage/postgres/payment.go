package postgres

import (
	"database/sql"
)

type PaymentRepo struct {
	Db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{Db: db}
}
