package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	pb "order_service/genproto/order"
	"order_service/models"
	"time"

	"github.com/google/uuid"
)

type OrderRepo struct {
	Db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{Db: db}
}

func (o *OrderRepo) CreateOrder(ctx context.Context, order *pb.ReqCreateOrder, total float64) (*pb.OrderInfo, error) {
	query := `
	INSERT INTO orders (
		id, user_id, kitchen_id, items, total_amount, status, delivery_address,
		delivery_time, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()
	deliveryTime := now.Add(time.Minute * 15).Format(time.RFC3339)
	createdAt := now.Format(time.RFC3339)
	updatedAt := now.Format(time.RFC3339)

	res := &pb.OrderInfo{
		Id:              uuid.NewString(),
		UserId:          order.UserId,
		KitchenId:       order.KitchenId,
		Items:           order.Items,
		TotalAmount:     total,
		Status:          "preparing",
		DeliveryAddress: order.DeliveryAddress,
		DeliveryTime:    deliveryTime,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	data, err := json.Marshal(res.Items)
	if err != nil {
		return nil, err
	}

	_, err = o.Db.ExecContext(ctx, query, res.Id, res.UserId, res.KitchenId, string(data), total, res.Status,
		res.DeliveryAddress, res.DeliveryTime, res.CreatedAt, res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *OrderRepo) UpdateOrderStatus(ctx context.Context, status *pb.Status) (*pb.StatusRes, error) {
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

func (o *OrderRepo) GetOrderById(ctx context.Context, id string) (*pb.OrderInfo, error) {
	query := `
	select
		id, user_id, kitchen_id, items, total_amount, status, delivery_address, delivery_time, created_at, updated_at
	from
		orders
	where
		id = $1
	`

	items := ""
	order := pb.OrderInfo{}
	row := o.Db.QueryRowContext(ctx, query, id)
	err := row.Scan(&order.Id, &order.UserId, &order.KitchenId, &items, &order.TotalAmount, &order.Status,
		&order.DeliveryAddress, &order.DeliveryTime, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	itemsObj := []*pb.Item{}
	err = json.Unmarshal([]byte(items), &itemsObj)
	if err != nil {
		return nil, err
	}
	order.Items = itemsObj

	return &order, nil
}

func (o *OrderRepo) GetOrdersForUser(ctx context.Context, filter *pb.Filter) (*pb.Orders, error) {
	query := `
	select
		id,
		user_id,
		status,
		total_amount,
		delivery_time
	from
		orders
	where
		user_id = $1 
	`
	query += fmt.Sprintf(" offset %d", (filter.Page-1)*filter.Limit)
	query += fmt.Sprintf(" limit %d", filter.Limit)

	rows, err := o.Db.QueryContext(ctx, query, filter.Id)
	if err != nil {
		return nil, err
	}

	orders := pb.Orders{}

	for rows.Next() {
		var order pb.OrderShortInfo

		err := rows.Scan(&order.Id, &order.UserId, &order.Status, &order.TotalAmount, &order.DeliveryTime)
		if err != nil {
			return nil, err
		}
		orders.Orders = append(orders.Orders, &order)
	}

	orders.Total = int64(o.GetOrderCount(ctx))
	orders.Limit = filter.Limit
	orders.Page = filter.Page

	return &orders, rows.Err()
}

func (o *OrderRepo) GetOrdersForChef(ctx context.Context, filter *pb.Filter) (*pb.Orders, error) {
	query := `
	select
		id,
		user_id,
		status,
		total_amount,
		delivery_time
	from
		orders
	where
		kitchen_id = $1 
	`
	query += fmt.Sprintf(" offset %d", (filter.Page-1)*filter.Limit)
	query += fmt.Sprintf(" limit %d", filter.Limit)

	rows, err := o.Db.QueryContext(ctx, query, filter.Id)
	if err != nil {
		return nil, err
	}

	orders := pb.Orders{}

	for rows.Next() {
		var order pb.OrderShortInfo

		err := rows.Scan(&order.Id, &order.UserId, &order.Status, &order.TotalAmount, &order.DeliveryTime)
		if err != nil {
			return nil, err
		}
		orders.Orders = append(orders.Orders, &order)
	}

	orders.Total = int64(o.GetOrderCount(ctx))
	orders.Limit = filter.Limit
	orders.Page = filter.Page

	return &orders, rows.Err()
}

func (o *OrderRepo) DeleteOrder(ctx context.Context, id string) error {
	query := `
	update
		orders
	set
		deleted_at = now()
	where
		id = $1 and deleted_at is null 
	`

	_, err := o.Db.ExecContext(ctx, query, id)

	return err
}

func (o *OrderRepo) ValidateOrderId(ctx context.Context, id string) error {
	query := `
	SELECT 
		1
	FROM 
		orders
	WHERE 
		id = $1
	`

	var exists int
	err := o.Db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order ID %s does not exist", id)
		}
		return fmt.Errorf("error checking order ID %s: %v", id, err)
	}

	return nil
}

func (o *OrderRepo) GetOrderCount(ctx context.Context) int {
	query := `
	select
		count(*)
	where
		deleted_at is null
	`
	count := 0
	err := o.Db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0
	}

	return count
}

func (o *OrderRepo) GetKitchenStatistics(ctx context.Context, filter *pb.DateFilter) (*pb.KitchenStatistics, error) {
	query := `
	with dish_data as (
		select
			jsonb_array_elements(items) as dish
		from
			orders
		where
			deleted_at is null and kitchen_id = $1 and
			created_at between $2 and $3
	),
	dish_stats as (
		select
			dish ->> 'dish_id' as id,
			count(*) as orders_count
		from
			dish_data
		group by
			dish ->> 'dish_id'
		order by
			orders_count desc
	)
	select
		id,
		orders_count
	from
		dish_stats
	`

	kitchenStats := pb.KitchenStatistics{}

	rows, err := o.Db.QueryContext(ctx, query, filter.Id, filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dishStats := []*pb.DishStats{}
	for rows.Next() {
		dishStat := &pb.DishStats{}
		err := rows.Scan(&dishStat.Id, &dishStat.OrdersCount)
		if err != nil {
			return nil, err
		}
		dishStats = append(dishStats, dishStat)
	}

	kitchenStats.TopDishes = dishStats

	return &kitchenStats, rows.Err()
}


func (o *OrderRepo) GetRevenueStatsForKitchen(ctx context.Context, filter *pb.DateFilter) (*models.RevenueStats, error){
	query := `
	select
		count(*),
		sum(total_amount)
	from
		orders
	where
		kitchen_id = $1 and
			created_at between $2 and $3
	`
	stats := models.RevenueStats{}
	row := o.Db.QueryRowContext(ctx, query, filter.Id, filter.StartDate, filter.EndDate)

	err := row.Scan(&stats.TotalOrders, &stats.Revenue)

	return &stats, err
}

func (o *OrderRepo) GetUserStatistics(ctx context.Context, filter *pb.DateFilter) (*pb.UserStatistics, error) {
	query := `
	with kitchen_data as (
		select
			kitchen_id,
			total_amount
		from
			orders
		where
			deleted_at is null and user_id = $1 and
			created_at between $2 and $3
	)
	select
		kitchen_id,
		count(*),
		sum(total_amount)
	from
		kitchen_data
	group by
		kitchen_id
	`

	userStats := pb.UserStatistics{}

	rows, err := o.Db.QueryContext(ctx, query, filter.Id, filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kitchenStats := []*pb.KitchenStats{}
	for rows.Next() {
		kitchenStat := &pb.KitchenStats{}
		err := rows.Scan(&kitchenStat.Id, &kitchenStat.OrdersCount, &kitchenStat.TotalSpent)
		if err != nil {
			return nil, err
		}
		kitchenStats = append(kitchenStats, kitchenStat)
	}

	userStats.FavoriteKitchens = kitchenStats

	return &userStats, rows.Err()
}
