package postgres

import (
	"context"
	"database/sql"
	"fmt"
	pb "order_service/genproto/payment"
	"time"

	"github.com/google/uuid"
)

type PaymentRepo struct {
	Db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{Db: db}
}

func (p *PaymentRepo) CreatePayment(ctx context.Context, req *pb.ReqCreatePayment, amount float64) (*pb.PaymentInfo, error) {
	query := `
	insert into
		payments(
		id,
		order_id,
		card_number,
		amount,
		status,
		payment_method,
		transaction_id,
		created_at,
		updated_at), 
	values($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	currentTime := time.Now().Format(time.RFC3339)
	res := pb.PaymentInfo{
		Id:            uuid.NewString(),
		OrderId:       req.OrderId,
		Amount:        amount,
		Status:        "Paid",
		TransactionId: "",
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
	}
	_, err := p.Db.ExecContext(ctx, query, res.Id, res.OrderId, req.CardNumber, amount, res.Status, req.PaymentMethod, res.TransactionId,
		res.CreatedAt, res.UpdatedAt)

	return &res, err
}

func (p *PaymentRepo) ValidateReviewId(ctx context.Context, id string) error {
	query := `
	SELECT 
		1
	FROM 
		payments
	WHERE 
		id = $1
	`

	var exists int
	err := p.Db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("review ID %s does not exist", id)
		}
		return fmt.Errorf("error checking review ID %s: %v", id, err)
	}

	return nil
}
