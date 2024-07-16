package main

import (
	"net"
	"order_service/config"
	pbd "order_service/genproto/dish"
	pbo "order_service/genproto/order"
	pbp "order_service/genproto/payment"
	pbr "order_service/genproto/review"
	"order_service/models"
	"order_service/pkg/logger"
	"order_service/service"
	"order_service/storage/postgres"
	"order_service/storage/redis"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	log, err := logger.New("debug", "development", cfg.LOG_PATH)
	if err != nil {
		panic(err)
	}

	postgresDb, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal("Cannot connect to Postgres", zap.Error(err))
		return
	}
	defer postgresDb.Close()

	redisDb, err := redis.ConnectDB()
	if err != nil {
		log.Fatal("Cannot connect to Redis", zap.Error(err))
		return
	}
	defer redisDb.Close()

	systemConfig := &models.SystemConfig{
		Config:     cfg,
		PostgresDb: postgresDb,
		RedisDb:    redisDb,
		Logger:     log,
	}

	listener, err := net.Listen("tcp", cfg.ORDER_SERVICE_PORT)
	if err != nil {
		systemConfig.Logger.Fatal("Failed to listen tcp", zap.Error(err))
		return
	}

	server := grpc.NewServer()

	pbd.RegisterDishServer(server, service.NewDishService(systemConfig))
	pbo.RegisterOrderServer(server, pbo.UnimplementedOrderServer{})
	pbp.RegisterPaymentServer(server, pbp.UnimplementedPaymentServer{})
	pbr.RegisterReviewServer(server, pbr.UnimplementedReviewServer{})

	err = server.Serve(listener)
	if err != nil {
		systemConfig.Logger.Fatal("grpc Failed to serve listener", zap.Error(err))
		return
	}
}
