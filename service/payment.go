package service

import (
	"order_service/models"
	"order_service/pkg/connections"
	"order_service/storage/postgres"

	pbk "order_service/genproto/kitchen"
	pb "order_service/genproto/payment"
	pbu "order_service/genproto/user"

	"go.uber.org/zap"
)

type PaymentService struct {
	paymentRepo   *postgres.PaymentRepo
	kitchenClient pbk.KitchenClient
	userClient    pbu.UserServiceClient
	log           *zap.Logger
	pb.UnimplementedPaymentServer
}

func NewPaymentService(sysConfig *models.SystemConfig) *PaymentService {
	return &PaymentService{
		paymentRepo:   postgres.NewPaymentRepo(sysConfig.PostgresDb),
		kitchenClient: connections.NewKitchenService(sysConfig),
		userClient:    connections.NewUserService(sysConfig),
		log:           sysConfig.Logger,
	}
}

