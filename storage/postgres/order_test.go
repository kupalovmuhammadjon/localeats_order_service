package postgres

import (
	"context"
	pb "order_service/genproto/order"
	"testing"
)

func newOrderepo() *OrderRepo {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	return &OrderRepo{Db: db}
}

func TestCreateOrder(t *testing.T) {
	o := newOrderepo()

	req := &pb.ReqCreateOrder{
		KitchenId:       "cdffffd7-67f2-4f96-b0b3-6d6b6bb85724",
		UserId:          "acdb0273-cb22-4168-9caf-360642cff29a",
		Items:           []*pb.Item{&pb.Item{
			DishId: "cdffffd7-67f2-4f96-b0b3-6d6b6bb85724",
			Quantity: 1,
		}},
		DeliveryAddress: "hgf",
		DeliveryTime:    "",
	}
	_, err := o.CreateOrder(context.Background(), req, 345)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOrderStatus(t *testing.T){

	o := newOrderepo()

	status := pb.Status{
		Id: "6e5c785c-d427-4ab5-afac-ed00400c08c7",
		Status: "delivering" ,
	}
	_, err := o.UpdateOrderStatus(context.Background(), &status)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOrderById(t *testing.T){

	o := newOrderepo()

	_, err := o.GetOrderById(context.Background(), "c2245570-9d3d-4a05-a8c2-1f6b8863eb64")
	if err != nil {
		t.Error(err)
	}
}

func TestGetOrdersForUser(t *testing.T){

	o := newOrderepo()

	filter := pb.Filter{
		Id:    "8529dbef-1313-4c78-b990-7a84ecb7d2c3",
		Page:  1,
		Limit: 10,
	}
	_, err := o.GetOrdersForUser(context.Background(), &filter)
	if err != nil {
		t.Error(err)
	}
}


func TestGetOrdersForChef(t *testing.T){

	o := newOrderepo()

	filter := pb.Filter{
		Id:    "8529dbef-1313-4c78-b990-7a84ecb7d2c3",
		Page:  1,
		Limit: 10,
	}
	_, err := o.GetOrdersForChef(context.Background(), &filter)
	if err != nil {
		t.Error(err)
	}
}

