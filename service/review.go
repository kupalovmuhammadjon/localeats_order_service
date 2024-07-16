package service

import (
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
