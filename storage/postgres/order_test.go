package postgres

import (
	"context"
	pb "order_service/genproto/order"
	"testing"
)

func newOrderepo() *Orderepo {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	return &Orderepo{Db: db}
}

func TestCreateOrder(t *testing.T) {
	o := newOrderepo()

	req := &pb.ReqCreateOrder{
		KitchenId:       "cdffffd7-67f2-4f96-b0b3-6d6b6bb85724",
		UserId:          "acdb0273-cb22-4168-9caf-360642cff29a",
		Items:           []*pb.Item{},
		DeliveryAddress: "hgf",
	}
	_, err := o.CreateOrder(context.Background(), req, 345)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOrderStatus(t *testing.T){

	o := newOrderepo()

	status := pb.Status{
		Id: "8529dbef-1313-4c78-b990-7a84ecb7d2c3",
		Status: "delivering" ,
	}
	_, err := o.UpdateOrderStatus(context.Background(), &status)
	if err != nil {
		t.Error(err)
	}
	
}

