package postgres

import (
	"context"
	pb "order_service/genproto/review"
	"testing"
)

func newReviewRepo() *ReviewRepo {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	return &ReviewRepo{Db: db}
}

func TestCreateReview(t *testing.T) {
	r := newReviewRepo()

	req := pb.ReqCreateReview{
		OrderId:   "",
		UserId:    "",
		KitchenId: "",
		Rating:    0,
		Comment:   "",
	}
	_, err := r.CreateReview(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
}