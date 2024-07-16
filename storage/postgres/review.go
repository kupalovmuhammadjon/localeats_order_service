package postgres

import (
	"database/sql"
)

type ReviewRepo struct {
	Db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{Db: db}
}
