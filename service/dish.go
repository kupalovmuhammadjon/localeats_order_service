package service

import (
	"context"
	"order_service/models"
	"order_service/pkg/connections"
	"order_service/storage/postgres"

	pb "order_service/genproto/dish"
	pbk "order_service/genproto/kitchen"
	pbu "order_service/genproto/user"

	"go.uber.org/zap"
)

type DishService struct {
	dishRepo      *postgres.DishRepo
	kitchenClient pbk.KitchenClient
	userClient    pbu.UserServiceClient
	log           *zap.Logger
	pb.UnimplementedDishServer
}

func NewDishService(sysConfig *models.SystemConfig) *DishService {
	return &DishService{
		dishRepo:      postgres.NewDishRepo(sysConfig.PostgresDb),
		kitchenClient: connections.NewKitchenService(sysConfig),
		userClient:    connections.NewUserService(sysConfig),
		log:           sysConfig.Logger,
	}
}

func (d *DishService) CreateDish(ctx context.Context, dish *pb.ReqCreateDish) (*pb.DishInfo, error) {
	_, err := d.kitchenClient.ValidateKitchenId(ctx, &pbk.Id{Id: dish.KitchenId})
	if err != nil {
		d.log.Info("Invalid kitchen Id ", zap.Error(err))
		return nil, err
	}
	res, err := d.dishRepo.CreateDish(ctx, dish)
	if err != nil {
		d.log.Error("failed to create dish ", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (d *DishService) UpdateDish(ctx context.Context, dish *pb.ReqUpdateDish) (*pb.DishInfo, error) {
	res, err := d.dishRepo.UpdateDish(ctx, dish)
	if err != nil {
		d.log.Error("failed to update dish ", zap.Error(err))
		return nil, err
	}

	kitchen, err := d.kitchenClient.GetKitchenById(ctx, &pbk.Id{Id: res.KitchenId})
	if err != nil {
		d.log.Error("failed to get kitchen by Id for dish ", zap.Error(err))
		return nil, err
	}
	res.KitchenName = kitchen.Name

	return res, nil
}

func (d *DishService) GetDishes(ctx context.Context, filter *pb.Pagination) (*pb.Dishes, error) {

	_, err := d.kitchenClient.ValidateKitchenId(ctx, &pbk.Id{Id: filter.Id})
	if err != nil {
		d.log.Info("Invalid kitchen Id ", zap.Error(err))
		return nil, err
	}

	res, err := d.dishRepo.GetDishes(ctx, filter)
	if err != nil {
		d.log.Error("failed to get dishes ", zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(res.Dishes); i++ {
		kitchen, err := d.kitchenClient.GetKitchenById(ctx, &pbk.Id{Id: res.Dishes[i].KitchenId})
		if err != nil {
			d.log.Error("failed to get kitchen by Id for dish ", zap.Error(err))
			return nil, err
		}
		res.Dishes[i].KitchenName = kitchen.Name
	}

	return res, nil
}

func (d *DishService) GetDishById(ctx context.Context, id *pb.Id) (*pb.DishInfo, error) {
	res, err := d.dishRepo.GetDishById(ctx, id)
	if err != nil {
		d.log.Error("failed to get dish by id ", zap.Error(err))
		return nil, err
	}

	kitchen, err := d.kitchenClient.GetKitchenById(ctx, &pbk.Id{Id: res.KitchenId})
	if err != nil {
		d.log.Error("failed to get kitchen by Id for dish ", zap.Error(err))
		return nil, err
	}
	res.KitchenName = kitchen.Name

	return res, nil
}

func (d *DishService) DeleteDish(ctx context.Context, id *pb.Id) (*pb.Void, error) {
	err := d.dishRepo.DeleteDish(ctx, id.Id)
	if err != nil {
		d.log.Error("failed to delete dish by id ", zap.Error(err))
		return nil, err
	}

	return &pb.Void{}, nil
}

func (d *DishService) ValidateDishId(ctx context.Context, id *pb.Id) (*pb.Void, error) {
	err := d.dishRepo.ValidateDishId(ctx, id.Id)
	if err != nil {
		d.log.Info("invalid dish id ", zap.Error(err))
		return nil, err
	}

	return &pb.Void{}, nil
}

func (d *DishService) UpdateNutritionInfo(ctx context.Context, info *pb.NutritionInfo) (*pb.DishInfo, error) {
	dish, err := d.dishRepo.UpdateNutritionInfo(ctx, info)
	if err != nil {
		d.log.Info("failed to update NutritionInfo ", zap.Error(err))
		return nil, err
	}

	return dish, nil
}

func (d *DishService) RecommendDishes(ctx context.Context, filter *pb.Filter) (*pb.Recommendations, error) {
	pref, err := d.userClient.GetUserPreference(ctx, &pbu.Id{Id: filter.Id})
	if err != nil {
		d.log.Error("failed to get user preferences for recommend dish ", zap.Error(err))
		return nil, err
	}

	ids, err := d.kitchenClient.GetKitchenIdsByCusineType(ctx, &pbk.Cusine{Cusine: pref.CuisineType})
	if err != nil {
		d.log.Error("failed to get kitchen ids ", zap.Error(err))
		return nil, err
	}
	res, err := d.dishRepo.RecommendDishes(ctx, filter, pref, ids.Ids)
	if err != nil {
		d.log.Error("failed to recommend dishes ", zap.Error(err))
		return nil, err
	}

	total, err := d.dishRepo.GetTotalRecommendation(ctx, filter, pref, ids.Ids)
	if err != nil {
		d.log.Error("failed to get total number of dishes ", zap.Error(err))
		return nil, err
	}
	res.Total = int32(total)
	res.Page = filter.Page
	res.Limit = filter.Limit

	return res, nil
}
