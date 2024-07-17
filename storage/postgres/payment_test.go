package postgres

import (
	"context"
	pb "order_service/genproto/payment"
	"testing"
)

func newPaymentRepo() *PaymentRepo {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	return &PaymentRepo{Db: db}
}

func TestCreatePayment(t *testing.T) {
	p := newPaymentRepo()

	req := pb.ReqCreatePayment{
		OrderId:       "8529dbef-1313-4c78-b990-7a84ecb7d2c3",
		PaymentMethod: "",
		CardNumber:    "",
		ExpiryDate:    "",
		Cvv:           "",
	}	

	_, err := p.CreatePayment(context.Background(), &req, 4)
	if err != nil {
		t.Error(err)
	}
}