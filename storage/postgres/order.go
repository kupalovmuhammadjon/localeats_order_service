package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	pb "order_service/genproto/order"
	"time"

	"github.com/google/uuid"
)

type Orderepo struct {
	Db *sql.DB
}

func NewOrderepo(db *sql.DB) *Orderepo {
	return &Orderepo{Db: db}
}

func (o *Orderepo) CreateOrder(ctx context.Context, order *pb.ReqCreateOrder, total float64) (*pb.OrderInfo, error) {
	query := `
	insert into
		orders(
			id, user_id, kitchen_id, items, total_amount, status, delivery_address,
			delivery_time,
			created_at,
			updated_at)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	res := &pb.OrderInfo{
		Id:              uuid.NewString(),
		UserId:          order.UserId,
		KitchenId:       order.KitchenId,
		Items:           order.Items,
		TotalAmount:     total,
		Status:          "preparing",
		DeliveryAddress: order.DeliveryAddress,
		DeliveryTime:    time.Now().Add(time.Minute*15).Format(time.RFC3339),
		CreatedAt:       time.Now().Format(time.RFC3339),
		UpdatedAt:       time.Now().Format(time.RFC3339),
	}
	data, err := json.Marshal(res.Items)
	if err != nil {
		return nil, err
	}

	_, err = o.Db.ExecContext(ctx, query, res.Id, res.UserId, res.KitchenId, string(data), total, res.Status,
		res.DeliveryAddress, res.DeliveryTime, res.CreatedAt, res.UpdatedAt)

	return res, err
}

func (o *Orderepo) UpdateOrderStatus(ctx context.Context, status *pb.Status) (*pb.StatusRes, error) {
	query := `
	update
		orders
	set
		status = $1,
		updated_at = $2
	where
		id = $3
	`

	res := &pb.StatusRes{
		Id:        status.Id,
		Status:    status.Status,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	_, err := o.Db.ExecContext(ctx, query, res.Status, res.UpdatedAt, res.Id)

	return res, err
}
