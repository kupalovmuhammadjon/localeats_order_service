package service

import (
	"context"
	"order_service/models"
	"order_service/pkg/connections"
	"order_service/storage/postgres"

	pbk "order_service/genproto/kitchen"
	pb "order_service/genproto/order"
	pbd "order_service/genproto/dish"
	pbu "order_service/genproto/user"

	"go.uber.org/zap"
)

type OrderService struct {
	orderRepo     *postgres.OrderRepo
	dishRepo     *postgres.DishRepo
	kitchenClient pbk.KitchenClient
	userClient    pbu.UserServiceClient
	log           *zap.Logger
	pb.UnimplementedOrderServer
}

func NewOrderService(sysConfig *models.SystemConfig) *OrderService {
	return &OrderService{
		orderRepo:     postgres.NewOrderRepo(sysConfig.PostgresDb),
		dishRepo: postgres.NewDishRepo(sysConfig.PostgresDb),
		kitchenClient: connections.NewKitchenService(sysConfig),
		userClient:    connections.NewUserService(sysConfig),
		log:           sysConfig.Logger,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, order *pb.ReqCreateOrder) (*pb.OrderInfo, error) {

	_, err := o.kitchenClient.ValidateKitchenId(ctx, &pbk.Id{Id: order.KitchenId})
	if err != nil {
		o.log.Info("invalid kitchen id ", zap.Error(err))
		return nil, err
	}
	_, err = o.userClient.ValidateUserId(ctx, &pbu.Id{Id: order.KitchenId})
	if err != nil {
		o.log.Info("invalid user id ", zap.Error(err))
		return nil, err
	}

	var total float64
	for _, item := range order.Items{
		dish, err := o.dishRepo.GetDishById(ctx, &pbd.Id{Id: item.DishId})
		if err != nil {
			o.log.Error("failed to get dish by id for order ", zap.Error(err))
			return nil, err
		}
		total += float64(dish.Price)
	}
	res, err := o.orderRepo.CreateOrder(ctx, order, total)
	if err != nil {
		o.log.Error("failed to create order ", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (o *OrderService) UpdateOrderStatus(ctx context.Context, status *pb.Status) (*pb.StatusRes, error)  {
	res, err := o.orderRepo.UpdateOrderStatus(ctx, status)
	if err != nil {
		o.log.Error("failed to update status of order ", zap.Error(err))
		return nil, err
	}

	return res, err
}

func (o *OrderService) GetOrderById(ctx context.Context, id *pb.Id) (*pb.OrderInfo, error)  {
	res, err := o.orderRepo.GetOrderById(ctx, id.Id)
	if err != nil {
		o.log.Error("failed to get order by id ", zap.Error(err))
		return nil, err
	}

	return res, err
}

func (o *OrderService) GetOrdersForUser(ctx context.Context, filter *pb.Filter) (*pb.Orders, error)  {
	res, err := o.orderRepo.GetOrdersForUser(ctx, filter)
	if err != nil {
		o.log.Error("failed to get orders for user ", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(res.Orders); i++{
		user, err := o.userClient.GetProfile(ctx, &pbu.Id{Id: res.Orders[i].UserId})
		if err != nil {
			o.log.Error("failed to get user profile for order ", zap.Error(err))
			return nil, err
		}
		res.Orders[i].Username = user.Username
	}

	return res, err
}

func (o *OrderService) GetOrdersForChef(ctx context.Context, filter *pb.Filter) (*pb.Orders, error)  {
	res, err := o.orderRepo.GetOrdersForChef(ctx, filter)
	if err != nil {
		o.log.Error("failed to get orders for chef ", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(res.Orders); i++{
		user, err := o.userClient.GetProfile(ctx, &pbu.Id{Id: res.Orders[i].UserId})
		if err != nil {
			o.log.Error("failed to get user profile for order ", zap.Error(err))
			return nil, err
		}
		res.Orders[i].Username = user.Username
	}

	return res, err
}

func (o *OrderService) DeleteOrder(ctx context.Context, id *pb.Id) (*pb.Void, error)  {
	err := o.orderRepo.DeleteOrder(ctx, id.Id)
	if err != nil {
		o.log.Error("failed to delete order ", zap.Error(err))
		return nil, err
	}

	return &pb.Void{}, err
}

func (o *OrderService) ValidateOrderId(ctx context.Context, id *pb.Id) (*pb.Void, error)  {
	err := o.orderRepo.ValidateOrderId(ctx, id.Id)
	if err != nil {
		o.log.Info("not valid order Id ", zap.Error(err))
		return nil, err
	}

	return &pb.Void{}, err
}