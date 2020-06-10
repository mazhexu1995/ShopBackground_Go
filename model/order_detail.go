package model

/**
 * 用户订单详情结构体
 */
type OrderDetail struct {
	UserOrder `xorm:"extends"`
	Address   `xorm:"extends"`
	Shop      `xorm:"extends"`
}
