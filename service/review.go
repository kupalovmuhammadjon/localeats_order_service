package service

import (
	"context"
	"order_service/models"
	"order_service/pkg/connections"
	"order_service/storage/postgres"

	pbk "order_service/genproto/kitchen"
	pb "order_service/genproto/review"
	pbu "order_service/genproto/user"

	"go.uber.org/zap"
)

type ReviewService struct {
	reviewRepo    *postgres.ReviewRepo
	kitchenClient pbk.KitchenClient
	userClient    pbu.UserServiceClient
	log           *zap.Logger
	pb.UnimplementedReviewServer
}

func NewReviewService(sysConfig *models.SystemConfig) *ReviewService {
	return &ReviewService{
		reviewRepo:    postgres.NewReviewRepo(sysConfig.PostgresDb),
		kitchenClient: connections.NewKitchenService(sysConfig),
		userClient:    connections.NewUserService(sysConfig),
		log:           sysConfig.Logger,
	}
}

func (r *ReviewService) CreateReview(ctx context.Context, req *pb.ReqCreateReview) (*pb.ReviewInfo, error){
	res, err := r.reviewRepo.CreateReview(ctx, req)
	if err != nil {
		r.log.Error("failed to create review ", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (r *ReviewService) GetReviewsByKitchenId(ctx context.Context, filter *pb.Filter) (*pb.Reviews, error){
	res, err := r.reviewRepo.GetReviewsByKitchenId(ctx, filter)
	if err != nil {
		r.log.Error("failed to create review ", zap.Error(err))
		return nil, err
	}
	
	return res, nil
}

func (r *ReviewService) DeleteComment(ctx context.Context, id *pb.Id) (*pb.Void, error){
	err := r.reviewRepo.DeleteReview(ctx, id.Id)
	if err != nil {
		r.log.Error("failed to create review ", zap.Error(err))
		return nil, err
	}
	
	return &pb.Void{}, err
}
