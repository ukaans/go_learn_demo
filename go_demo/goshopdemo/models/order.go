package models

type Order struct {
	Id          int
	OrderId     string
	Uid         int
	AllPrice    float64
	Phone       string
	Name        string
	Address     string
	PayStatus   int // 支付状态： 0 表示未支付     1 已经支付
	PayType     int // 支付类型： 0 alipay    1 wechat  2 信用点
	OrderStatus int // 订单状态： 0 已下单  1 已付款  2 已配货  3、发货   4、交易成功   5、退货   6、取消
	AddTime     int //下单时间
	PayTime     int //支付时间
	CancelTime  int //取消时间

	OrderItem []OrderItem `gorm:"foreignKey:OrderId;references:Id"`
}

func (Order) TableName() string {
	return "order"
}
