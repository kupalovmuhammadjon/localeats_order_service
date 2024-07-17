package service

import (
	"context"
	"order_service/models"
	"order_service/pkg/connections"
	"order_service/storage/postgres"

	pbd "order_service/genproto/dish"
	pbk "order_service/genproto/kitchen"
	pb "order_service/genproto/order"
	pbu "order_service/genproto/user"

	"go.uber.org/zap"
)

type OrderService struct {
	orderRepo     *postgres.OrderRepo
	dishRepo      *postgres.DishRepo
	reviewRepo    *postgres.ReviewRepo
	kitchenClient pbk.KitchenClient
	userClient    pbu.UserServiceClient
	log           *zap.Logger
	pb.UnimplementedOrderServer
}

func NewOrderService(sysConfig *models.SystemConfig) *OrderService {
	return &OrderService{
		orderRepo:     postgres.NewOrderRepo(sysConfig.PostgresDb),
		dishRepo:      postgres.NewDishRepo(sysConfig.PostgresDb),
		reviewRepo:    postgres.NewReviewRepo(sysConfig.PostgresDb),
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
	_, err = o.userClient.ValidateUserId(ctx, &pbu.Id{Id: order.UserId})
	if err != nil {
		o.log.Info("invalid user id ", zap.Error(err))
		return nil, err
	}

	var total float64
	for _, item := range order.Items {
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

func (o *OrderService) UpdateOrderStatus(ctx context.Context, status *pb.Status) (*pb.StatusRes, error) {
	res, err := o.orderRepo.UpdateOrderStatus(ctx, status)
	if err != nil {
		o.log.Error("failed to update status of order ", zap.Error(err))
		return nil, err
	}

	return res, err
}

func (o *OrderService) GetOrderById(ctx context.Context, id *pb.Id) (*pb.OrderInfo, error) {
	res, err := o.orderRepo.GetOrderById(ctx, id.Id)
	if err != nil {
		o.log.Error("failed to get order by id ", zap.Error(err))
		return nil, err
	}

	return res, err
}

func (o *OrderService) GetOrdersForUser(ctx context.Context, filter *pb.Filter) (*pb.Orders, error) {
	res, err := o.orderRepo.GetOrdersForUser(ctx, filter)
	if err != nil {
		o.log.Error("failed to get orders for user ", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(res.Orders); i++ {
		user, err := o.userClient.GetProfile(ctx, &pbu.Id{Id: res.Orders[i].UserId})
		if err != nil {
			o.log.Error("failed to get user profile for order ", zap.Error(err))
			return nil, err
		}
		res.Orders[i].Username = user.Username
	}

	return res, err
}

func (o *OrderService) GetOrdersForChef(ctx context.Context, filter *pb.Filter) (*pb.Orders, error) {
	res, err := o.orderRepo.GetOrdersForChef(ctx, filter)
	if err != nil {
		o.log.Error("failed to get orders for chef ", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(res.Orders); i++ {
		user, err := o.userClient.GetProfile(ctx, &pbu.Id{Id: res.Orders[i].UserId})
		if err != nil {
			o.log.Error("failed to get user profile for order ", zap.Error(err))
			return nil, err
		}
		res.Orders[i].Username = user.Username
	}

	return res, err
}

func (o *OrderService) DeleteOrder(ctx context.Context, id *pb.Id) (*pb.Void, error) {
	err := o.orderRepo.DeleteOrder(ctx, id.Id)
	if err != nil {
		o.log.Error("failed to delete order ", zap.Error(err))
		return nil, err
	}

	return &pb.Void{}, err
}

func (o *OrderService) ValidateOrderId(ctx context.Context, id *pb.Id) (*pb.Void, error) {
	err := o.orderRepo.ValidateOrderId(ctx, id.Id)
	if err != nil {
		o.log.Info("not valid order Id ", zap.Error(err))
		return nil, err
	}

	return &pb.Void{}, err
}

func (o *OrderService) GetKitchenStatistics(ctx context.Context, filter *pb.DateFilter) (*pb.KitchenStatistics, error) {
	reviewStats, err := o.reviewRepo.GetStatisticsOfReviews(ctx, filter.Id)
	if err != nil {
		o.log.Error("failed to get review stats", zap.Error(err))
		return nil, err
	}
	rev, err := o.orderRepo.GetRevenueStatsForKitchen(ctx, filter)
	if err != nil {
		o.log.Error("failed to get revenu stats", zap.Error(err))
		return nil, err
	}

	statistics, err := o.orderRepo.GetKitchenStatistics(ctx, filter)
	if err != nil {
		o.log.Error("failed to get kitchen stats", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(statistics.TopDishes); i++{
		dish, err := o.dishRepo.GetDishById(ctx, &pbd.Id{Id: statistics.TopDishes[i].Id})
		if err != nil {
			o.log.Info("failed to get dish by id", zap.Error(err))
		}
		statistics.TopDishes[i].Name = dish.Name
		statistics.TopDishes[i].Revenue = dish.Price * float32(statistics.TopDishes[i].OrdersCount)
	}

	statistics.AverageRating = reviewStats.AvarageRating
	statistics.TotalRevenue = rev.Revenue
	statistics.TotalOrders = int64(rev.TotalOrders)

	return statistics, nil
}

func (o *OrderService) GetUserStatistics(ctx context.Context, filter *pb.DateFilter) (*pb.UserStatistics, error) {

	statistics, err := o.orderRepo.GetUserStatistics(ctx, filter)
	if err != nil {
		o.log.Error("failed to get user stats", zap.Error(err))
		return nil, err
	}
	totalSpent := 0
	totalOrders := 0
	for i := 0; i < len(statistics.FavoriteKitchens); i++{
		kitchen, err := o.kitchenClient.GetKitchenById(ctx, &pbk.Id{Id: statistics.FavoriteKitchens[i].Id})
		if err != nil {
			o.log.Info("failed to get kitchen by id", zap.Error(err))
		}
		statistics.FavoriteKitchens[i].Name = kitchen.Name	
		totalSpent += int(statistics.FavoriteKitchens[i].TotalSpent)
		totalOrders += int(statistics.FavoriteKitchens[i].OrdersCount)
	}

	statistics.TotalSpent = float64(totalSpent)
	statistics.TotalOrders = int64(totalOrders)

	return statistics, nil
}
