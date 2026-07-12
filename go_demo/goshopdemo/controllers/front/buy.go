package front

import (
	"goshopdemo/models"

	"github.com/gin-gonic/gin"
)

type BuyController struct {
	BaseController
}

func (con BuyController) Checkout(c *gin.Context) {
	//1、获取购物车中选择的商品

	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)

	orderList := []models.Cart{}
	var allPrice float64
	var allNum int

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
			allNum += cartList[i].Num
		}
	}

	//2、获取当前用户的收货地址

	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	addressList := []models.Address{}
	models.DB.Where("uid = ?", user.Id).Order("id desc").Find(&addressList)

	con.Render(c, "itying/buy/checkout.html", gin.H{
		"orderList":   orderList,
		"allPrice":    allPrice,
		"allNum":      allNum,
		"addressList": addressList,
	})

}

/*
提交订单执行结算

	1、获取用户信息 获取用户的收货地址信息
	2、获取购买商品的信息
	3、把订单信息放在订单表，把商品信息放在商品表
	4、删除购物车里面的选中数据
	5、跳转到支付页面
*/
func (con BuyController) DoCheckout(c *gin.Context) {
	// 1、获取用户信息 获取用户的收货地址信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	addressResult := models.Address{}
	models.DB.Where("uid = ? AND default_address=1", user.Id).Find(&addressResult)

	// 2、获取购买商品的信息
	cartList := []models.Cart{}
	models.Cookie.Get(c, "cartList", &cartList)
	orderList := []models.Cart{}
	var allPrice float64
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
		}
	}
	// 3、把订单信息放在订单表，把商品信息放在商品表
	order := models.Order{
		OrderId:     models.GetOrderId(),
		Uid:         user.Id,
		AllPrice:    allPrice,
		Phone:       addressResult.Phone,
		Name:        addressResult.Name,
		Address:     addressResult.Address,
		PayStatus:   0,
		PayType:     0,
		OrderStatus: 0,
		AddTime:     int(models.GetUnix()),
	}

	err := models.DB.Create(&order).Error
	//增加数据成功以后可以通过  order.Id
	if err == nil {
		// 把商品信息放在商品对应的订单表
		for i := 0; i < len(orderList); i++ {
			orderItem := models.OrderItem{
				OrderId:      order.Id,
				Uid:          user.Id,
				ProductTitle: orderList[i].Title,
				ProductId:    orderList[i].Id,
				ProductImg:   orderList[i].GoodsImg,
				ProductPrice: orderList[i].Price,
				ProductNum:   orderList[i].Num,
				GoodsVersion: orderList[i].GoodsVersion,
				GoodsColor:   orderList[i].GoodsColor,
			}
			models.DB.Create(&orderItem)
		}
	}

	// 4、删除购物车里面的选中数据
	noSelectCartList := []models.Cart{}
	for i := 0; i < len(cartList); i++ {
		if !cartList[i].Checked {
			noSelectCartList = append(noSelectCartList, cartList[i])
		}
	}
	models.Cookie.Set(c, "cartList", noSelectCartList)

	c.Redirect(302, "/buy/pay")
}

// 支付
func (con BuyController) Pay(c *gin.Context) {
	c.String(200, "支付页面")
}
