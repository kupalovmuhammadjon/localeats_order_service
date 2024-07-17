package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	pb "order_service/genproto/dish"
	pbu "order_service/genproto/user"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DishRepo struct {
	Db *sql.DB
}

func NewDishRepo(db *sql.DB) *DishRepo {
	return &DishRepo{Db: db}
}

func (d *DishRepo) CreateDish(ctx context.Context, dish *pb.ReqCreateDish) (*pb.DishInfo, error) {
	query := `
	insert into
		dishes(
			id,
			kitchen_id,
			name,
			description,
			price,
			category,
			ingredients,
			available,
			created_at,
			updated_at
			)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	res := &pb.DishInfo{
		Id:            uuid.NewString(),
		KitchenId:     dish.KitchenId,
		Name:          dish.Name,
		Price:         dish.Price,
		Category:      dish.Category,
		Ingredients:   dish.Ingredients,
		Description:   dish.Description,
		Available:     dish.Available,
		Allergens:     []string{},
		NutritionInfo: "",
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	_, err := d.Db.ExecContext(ctx, query, res.Id, res.KitchenId, res.Name, res.Description, res.Price, res.Category,
		pq.Array(res.Ingredients), res.Available, res.CreatedAt, res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DishRepo) UpdateDish(ctx context.Context, dish *pb.ReqUpdateDish) (*pb.DishInfo, error) {
	query := `
	update
		dishes
	set
		name = $1,
		description = $2,
		price = $3,
		category = $4,
		ingredients = $5,
		available = $6,
		updated_at = now()
	where
		id = $7 and deleted_at is null
	returning id, kitchen_id, name, description, price, category, ingredients, allergens, nutrition_info, dietary_info, available,
	created_at, updated_at
	`

	res := &pb.DishInfo{}

	row := d.Db.QueryRowContext(ctx, query, dish.Name, dish.Description, dish.Price, dish.Category, pq.Array(dish.Ingredients),
		dish.Available, dish.Id)

	var nutritionInfo sql.NullString
	err := row.Scan(&res.Id, &res.KitchenId, &res.Name, &res.Description, &res.Price, &res.Category, pq.Array(&res.Ingredients),
		pq.Array(&res.Allergens), &nutritionInfo, pq.Array(&res.DietaryInfo), &res.Available, &res.CreatedAt, &res.UpdatedAt)

	res.NutritionInfo = nutritionInfo.String

	return res, err
}

func (d *DishRepo) GetDishes(ctx context.Context, pagination *pb.Pagination) (*pb.Dishes, error) {
	query := `
	select
		id, kitchen_id, price, category, available
	from
		dishes
	where
		deleted_at is null and kitchen_id = $1 and available = true
	`
	query += fmt.Sprintf(" offset %d", (pagination.Page-1)*pagination.Limit)
	query += fmt.Sprintf(" limit %d", pagination.Limit)

	rows, err := d.Db.QueryContext(ctx, query, pagination.Id)
	if err != nil {
		return nil, err
	}

	dishes := pb.Dishes{}
	for rows.Next() {
		dish := pb.DishShortInfo{}
		err := rows.Scan(&dish.Id, &dish.KitchenId, &dish.Price, &dish.Category, &dish.Available)
		if err != nil {
			return nil, err
		}
		dishes.Dishes = append(dishes.Dishes, &dish)
	}

	return &dishes, rows.Err()
}

func (d *DishRepo) GetDishById(ctx context.Context, id *pb.Id) (*pb.DishInfo, error) {
	query := `
	select
		id, kitchen_id, name, description, price, category, ingredients, allergens, nutrition_info, dietary_info, available,
	created_at, updated_at
	from
		dishes
	where
		deleted_at is null and id = $1
	`

	row := d.Db.QueryRowContext(ctx, query, id.Id)

	dish := &pb.DishInfo{}

	var nutritionInfo sql.NullString
	err := row.Scan(&dish.Id, &dish.KitchenId, &dish.Name, &dish.Description, &dish.Price,
		&dish.Category, pq.Array(&dish.Ingredients), pq.Array(&dish.Allergens), &nutritionInfo, pq.Array(&dish.DietaryInfo),
		&dish.Available, &dish.CreatedAt, &dish.UpdatedAt)
	if err != nil {
		return nil, err
	}
	dish.NutritionInfo = nutritionInfo.String

	return dish, row.Err()
}

func (d *DishRepo) DeleteDish(ctx context.Context, id string) error {
	query := `
	update
		dishes
	set
		deleted_at = now()
	where
		id = $1 and deleted_at is null 
	`

	_, err := d.Db.ExecContext(ctx, query, id)

	return err
}

func (d *DishRepo) ValidateDishId(ctx context.Context, id string) error {
	query := `
	SELECT 
		1
	FROM 
		dishes
	WHERE 
		id = $1
	`

	var exists int
	err := d.Db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("dish ID %s does not exist", id)
		}
		return fmt.Errorf("error checking dish ID %s: %v", id, err)
	}

	return nil
}

func (d *DishRepo) UpdateNutritionInfo(ctx context.Context, info *pb.NutritionInfo) (*pb.DishInfo, error) {
	query := `
	update
		dishes
	set
		allergens = $1,
		nutrition_info = $2,
		dietary_info = $3,
		updated_at = now()
	where
		id = $4 and deleted_at is null
	returning id, kitchen_id, name, description, price, category, ingredients, allergens, 
	nutrition_info, dietary_info, available, created_at, updated_at
	`

	res := &pb.DishInfo{}
	nutritions := map[string]int32{
		"calories":      info.Calories,
		"protein":       info.Protein,
		"carbohydrates": info.Carbohydrates,
		"fat":           info.Fat,
	}
	data, err := json.Marshal(nutritions)
	if err != nil {
		return nil, err
	}
	row := d.Db.QueryRowContext(ctx, query, pq.Array(info.Allergens), string(data), pq.Array(info.DietaryInfo), info.Id)

	var nutritionInfo sql.NullString
	err = row.Scan(&res.Id, &res.KitchenId, &res.Name, &res.Description, &res.Price, &res.Category, pq.Array(&res.Ingredients),
		pq.Array(&res.Allergens), &nutritionInfo, pq.Array(&res.DietaryInfo), &res.Available, &res.CreatedAt, &res.UpdatedAt)

	res.NutritionInfo = nutritionInfo.String

	return res, err
}

func (d *DishRepo) RecommendDishes(ctx context.Context, filter *pb.Filter, user *pbu.PreferencesRes, kitchens []string) (*pb.Recommendations, error) {
	/*
		Recommendation criterias
		user preferences (favourite kitchens, dieatery and cusine type)
		latest highly rated food
	*/

	query := `
	select
		id, kitchen_id, price, category, available
	from
		dishes
	where
		deleted_at is null and available = true
		and(
			kitchen_id = any($1) or 
			to_tsvector(dietary_info::text) @@ plainto_tsquery($2)
		)
	`

	query += fmt.Sprintf(" offset %d", (filter.Page-1)*filter.Limit)
	query += fmt.Sprintf(" limit %d", filter.Limit)

	search := strings.Join(user.DietaryPreferences, " ")

	rows, err := d.Db.QueryContext(ctx, query, pq.Array(kitchens), search)
	if err != nil {
		return nil, err
	}

	dishes := pb.Recommendations{}
	for rows.Next() {
		dish := pb.DishShortInfo{}
		err := rows.Scan(&dish.Id, &dish.KitchenId, &dish.Price, &dish.Category, &dish.Available)
		if err != nil {
			return nil, err
		}
		dishes.Dishes = append(dishes.Dishes, &dish)
	}

	return &dishes, rows.Err()
}

func (d *DishRepo) GetTotalRecommendation(ctx context.Context, filter *pb.Filter, user *pbu.PreferencesRes, kitchens []string) (int, error) {
	query := `
	select
		count(*)
	from
		dishes
	where
		deleted_at is null and available = true
		and(
			kitchen_id = any($1) or 
			to_tsvector(dietary_info::text) @@ plainto_tsquery($2)
		)
	`

	search := strings.Join(user.DietaryPreferences, " ")

	var total int
	err := d.Db.QueryRowContext(ctx, query, pq.Array(kitchens), search).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return total, nil
}
