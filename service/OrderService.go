package service

import (
	"ShopBackground/model"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
)

/**
 * 订单服务接口
 */
type OrderService interface {
	GetCount() (int64, error)
	GetOrderList(offset, limit int) []model.OrderDetail
}

/**
 * 订单服务
 */
type orderService struct {
	Engine *xorm.Engine
}

/**
 * 实例化OrderService服务对象
 */
func NewOrderService(db *xorm.Engine) OrderService {
	return &orderService{Engine: db}
}

/**
 * 获取订单列表
 */
func (orderService *orderService) GetOrderList(offset, limit int) []model.OrderDetail {

	orderList := make([]model.OrderDetail, 0)
	err := orderService.Engine.Table("user_order").
		Join("INNER", "order_status", " order_status.id = user_order.order_status_id ").
		Join("INNER", "user", " user.id = user_order.user_id").
		Join("INNER", "shop", " shop.shop_id = user_order.shop_id ").
		Join("INNER", "address", " address.address_id = user_order.address_id ").
		Find(&orderList)
	if err != nil {
		iris.New().Logger().Error(err.Error())
		panic(err.Error())
		return nil
	}
	iris.New().Logger().Info(orderList)
	return orderList
}

/**
 * 获取订单总数量
 */
func (orderService *orderService) GetCount() (int64, error) {
	count, err := orderService.Engine.Where(" del_flag = 0 ").Count(new(model.UserOrder))
	if err != nil {
		return 0, err
	}
	return count, nil
}

