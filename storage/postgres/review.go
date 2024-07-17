package postgres

import (
	"context"
	"database/sql"
	"fmt"
	pb "order_service/genproto/review"
	"order_service/models"
	"time"

	"github.com/google/uuid"
)

type ReviewRepo struct {
	Db *sql.DB
}

func NewReviewRepo(db *sql.DB) *ReviewRepo {
	return &ReviewRepo{Db: db}
}

func (r *ReviewRepo) CreateReview(ctx context.Context, review *pb.ReqCreateReview) (*pb.ReviewInfo, error) {
	query := `
	insert into
		reviews(
		id,
		order_id,
		user_id,
		kitchen_id,
		rating,
		comment,
		created_at,
		updated_at)
	values($1, $2, $3, $4, $5, $6, $7, $8)
	`

	currentTime := time.Now().Format(time.RFC3339)
	res := pb.ReviewInfo{
		Id:        uuid.NewString(),
		OrderId:   review.OrderId,
		UserId:    review.UserId,
		KitchenId: review.KitchenId,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	_, err := r.Db.ExecContext(ctx, query, res.Id, res.OrderId, res.UserId, res.KitchenId, res.Rating, res.Comment,
		res.CreatedAt, res.UpdatedAt)

	return &res, err
}

func (r *ReviewRepo) GetReviewsByKitchenId(ctx context.Context, filter *pb.Filter) (*pb.Reviews, error) {
	query := `
	select
		id,
		order_id,
		user_id,
		rating,
		comment,
		created_at,
		updated_at
	from
		reviews
	where
		kitchen_id = $1
	`

	reviews := pb.Reviews{}

	rows, err := r.Db.QueryContext(ctx, query, filter.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var review pb.ReviewShortInfo

		err = rows.Scan(&review.Id, &review.OrderId, &review.UserId, &review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reviews.Reviews = append(reviews.Reviews, &review)
	}
	stats, err := r.GetStatisticsOfReviews(ctx, filter.Id)
	if err != nil {
		return nil, err
	}
	reviews.Total = int64(stats.TotalNumberOfComments)
	reviews.AverageRating = stats.AvarageRating
	reviews.Page = filter.Page
	reviews.Limit = filter.Limit

	return &reviews, rows.Err()
}

func (r *ReviewRepo) GetStatisticsOfReviews(ctx context.Context, id string) (*models.ReviewsStats, error) {
	query := `
	select 
		count(*),
		round(avg(rating)::numeric, 2)
	from
		reviews
	where
		kitchen_id = $1
	`

	res := models.ReviewsStats{}

	err := r.Db.QueryRowContext(ctx, query, id).Scan(&res.TotalNumberOfComments, &res.AvarageRating)

	return &res, err
}

func (r *ReviewRepo) DeleteReview(ctx context.Context, id string) error {
	query := `
	update
		reviews
	set
		deleted_at = now()
	where
		id = $1 and deleted_at is null 
	`

	_, err := r.Db.ExecContext(ctx, query, id)

	return err
}

func (r *ReviewRepo) ValidateReviewId(ctx context.Context, id string) error {
	query := `
	SELECT 
		1
	FROM 
		reviews
	WHERE 
		id = $1
	`

	var exists int
	err := r.Db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("review ID %s does not exist", id)
		}
		return fmt.Errorf("error checking review ID %s: %v", id, err)
	}

	return nil
}

