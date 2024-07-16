package service

import (
	"context"
	"order_service/models"
	"order_service/storage/postgres"

	pb "order_service/genproto/dish"

	"go.uber.org/zap"
)

type DishService struct {
	dishRepo *postgres.DishRepo
	log      *zap.Logger
	pb.UnimplementedDishServer
}

func NewDishService(sysConfig *models.SystemConfig) *DishService {
	return &DishService{
		dishRepo: postgres.NewDishRepo(sysConfig.PostgresDb),
		log:      sysConfig.Logger,
	}
}

func (d *DishService) CreateDish(ctx context.Context, dish *pb.ReqCreateDish) (*pb.DishInfo, error) {
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

	return res, nil
}

func (d *DishService) GetDishes(ctx context.Context, filter *pb.Pagination) (*pb.Dishes, error) {
	res, err := d.dishRepo.GetDishes(ctx, filter)
	if err != nil {
		d.log.Error("failed to get dishes ", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (d *DishService) GetDishById(ctx context.Context, id *pb.Id) (*pb.DishInfo, error) {
	res, err := d.dishRepo.GetDishById(ctx, id)
	if err != nil {
		d.log.Error("failed to get dish by id ", zap.Error(err))
		return nil, err
	}

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
		d.log.Info("invalid dish id ", zap.Error(err))
		return nil, err
	}

	return dish, nil
}

