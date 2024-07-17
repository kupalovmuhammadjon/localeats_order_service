package service

import (
	"context"
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
	orderRepo     *postgres.OrderRepo
	kitchenClient pbk.KitchenClient
	userClient    pbu.UserServiceClient
	log           *zap.Logger
	pb.UnimplementedPaymentServer
}

func NewPaymentService(sysConfig *models.SystemConfig) *PaymentService {
	return &PaymentService{
		paymentRepo:   postgres.NewPaymentRepo(sysConfig.PostgresDb),
		orderRepo:     postgres.NewOrderRepo(sysConfig.PostgresDb),
		kitchenClient: connections.NewKitchenService(sysConfig),
		userClient:    connections.NewUserService(sysConfig),
		log:           sysConfig.Logger,
	}
}

func (p *PaymentService) CreatePayment(ctx context.Context, req *pb.ReqCreatePayment) (*pb.PaymentInfo, error) {
	order, err := p.orderRepo.GetOrderById(ctx, req.OrderId)
	if err != nil {
		p.log.Error("Failed to get order by id for id ", zap.Error(err))
		return nil, err
	}
	
	res, err := p.paymentRepo.CreatePayment(ctx, req, order.TotalAmount)
	if err != nil {
		p.log.Error("Failed to create payment ", zap.Error(err))
		return nil, err
	}

	return res, nil
}
