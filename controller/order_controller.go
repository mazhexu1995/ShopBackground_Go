package controller

import (
	"context"
	"ShopBackground/service"
)

type OrderController struct {
	Ctx context.Context
	Service service.OrderService
}
