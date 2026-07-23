package front

import (
	"fmt"
	"goshopdemo/models"

	"github.com/gin-contrib/sessions"
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

	//3、生成签名
	orderSign := models.Md5(models.GetRandomNum())
	session := sessions.Default(c)
	session.Set("orderSign", orderSign)
	session.Save()

	//4、判断orderList数据是否存在
	if len(orderList) == 0 {
		c.Redirect(302, "/")
		return
	}

	con.Render(c, "itying/buy/checkout.html", gin.H{
		"orderList":   orderList,
		"allPrice":    allPrice,
		"allNum":      allNum,
		"addressList": addressList,
		"orderSign":   orderSign,
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
	//0、防止重复提交订单
	orderSignClient := c.PostForm("orderSign")
	session := sessions.Default(c)
	orderSignSession := session.Get("orderSign")

	fmt.Println("=== DoCheckout Debug ===")
	fmt.Println("Client orderSign:", orderSignClient)
	fmt.Println("Session orderSign:", orderSignSession)

	if orderSignSession == nil {
		fmt.Println("错误：Session 中没有 orderSign")
		c.Redirect(302, "/")
		return
	}

	orderSignServer, ok := orderSignSession.(string)
	if !ok || orderSignClient == "" || orderSignClient != orderSignServer {
		fmt.Println("错误：orderSign 不匹配！")
		c.Redirect(302, "/")
		return
	}
	session.Delete("orderSign")
	session.Save()

	// 1、获取用户信息
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

	if len(orderList) == 0 {
		c.Redirect(302, "/cart")
		return
	}

	// 3、创建订单
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
	if err == nil {
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
	} else {
		fmt.Println("创建订单失败:", err)
		c.Redirect(302, "/cart")
		return
	}

	// 4、清理购物车（只保留未选中的商品，并清除Checked状态）
	noSelectCartList := []models.Cart{}
	for i := 0; i < len(cartList); i++ {
		if !cartList[i].Checked {
			cartList[i].Checked = false
			noSelectCartList = append(noSelectCartList, cartList[i])
		}
	}
	models.Cookie.Set(c, "cartList", noSelectCartList)

	// 5、重要：跳转到确认页（或支付页）
	c.Redirect(302, "/buy/pay?orderId="+models.String(order.Id))
}

// 订单提交成功确认页（无需orderSign验证）
// 订单确认页（去支付入口）
func (con BuyController) Confirm(c *gin.Context) {
	fmt.Println("=== Confirm Handler Called ===")
	fmt.Println("orderId:", c.Query("orderId"))
	orderId, err := models.Int(c.Query("orderId"))
	if err != nil || orderId <= 0 {
		// 兼容旧的 ?id= 参数
		orderId, _ = models.Int(c.Query("id"))
		if orderId <= 0 {
			c.Redirect(302, "/user/order")
			return
		}
	}

	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	if user.Id == 0 {
		c.Redirect(302, "/")
		return
	}

	order := models.Order{}
	models.DB.Where("id = ? AND uid = ?", orderId, user.Id).First(&order)

	if order.Id == 0 {
		c.Redirect(302, "/user/order?msg=order_not_found")
		return
	}

	// 查询订单商品
	orderItems := []models.OrderItem{}
	models.DB.Where("order_id = ?", orderId).Find(&orderItems)

	con.Render(c, "itying/buy/confirm.html", gin.H{
		"order":      order,
		"orderItems": orderItems,
	})
}

// 支付
func (con BuyController) Pay(c *gin.Context) {

	orderId, err := models.Int(c.Query("orderId"))
	if err != nil {
		c.Redirect(302, "/")
	}
	//获取用户信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	//获取订单信息
	order := models.Order{}
	models.DB.Where("id = ?", orderId).Find(&order)
	if order.Uid != user.Id {
		c.Redirect(302, "/")
		return
	}
	//获取订单对应的商品

	orderItems := []models.OrderItem{}
	models.DB.Where("order_id = ?", orderId).Find(&orderItems)

	con.Render(c, "itying/buy/pay.html", gin.H{
		"order":      order,
		"orderItems": orderItems,
		"allPrice":   order.AllPrice,
	})
}

// 信用点支付成功处理
func (con BuyController) PaySuccess(c *gin.Context) {
	orderId, err := models.Int(c.Query("orderId"))
	payType, _ := models.Int(c.Query("payType")) // 2 表示信用点

	if err != nil || orderId <= 0 {
		c.Redirect(302, "/")
		return
	}

	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	// 查询订单
	order := models.Order{}
	models.DB.Where("id = ? AND uid = ?", orderId, user.Id).First(&order)

	if order.Id == 0 {
		c.Redirect(302, "/")
		return
	}

	// 更新订单状态
	now := int(models.GetUnix())
	updateData := map[string]interface{}{
		"pay_status":   1,
		"pay_type":     payType, // 2 = 信用点
		"order_status": 1,
		"pay_time":     now,
	}

	err = models.DB.Model(&order).Updates(updateData).Error
	if err != nil {
		fmt.Println("支付成功更新失败:", err)
		return
	}

	// 显示支付成功提示页
	con.Render(c, "itying/buy/pay_success.html", gin.H{
		"order": order,
	})
}

func (con BuyController) Ship(c *gin.Context) {
	orderId, err := models.Int(c.Query("orderId"))
	fmt.Println("=== Ship Handler Called ===")
	fmt.Println("orderId from URL:", orderId)

	if err != nil || orderId <= 0 {
		fmt.Println("错误：orderId 无效")
		c.Redirect(302, "/user/order?msg=invalid_order")
		return
	}

	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	fmt.Println("当前登录用户ID:", user.Id)

	if user.Id == 0 {
		fmt.Println("错误：用户未登录")
		c.Redirect(302, "/")
		return
	}

	order := models.Order{}
	models.DB.Where("id = ? AND uid = ?", orderId, user.Id).First(&order)

	fmt.Println("查询到的订单ID:", order.Id)
	fmt.Println("订单当前 OrderStatus:", order.OrderStatus)
	fmt.Println("订单所属用户ID:", order.Uid)

	if order.Id == 0 {
		fmt.Println("错误：订单不存在或不属于当前用户")
		c.Redirect(302, "/user/order?msg=order_not_found")
		return
	}

	// 执行更新
	updateData := map[string]interface{}{
		"order_status":      3,
		"distribution_time": int(models.GetUnix()),
	}

	result := models.DB.Model(&order).Updates(updateData)
	fmt.Println("更新影响行数:", result.RowsAffected)
	fmt.Println("更新错误:", result.Error)

	if result.Error != nil {
		fmt.Println("数据库更新失败:", result.Error)
		c.Redirect(302, "/user/order?msg=ship_failed")
		return
	}

	fmt.Println("发货成功！订单ID:", orderId)
	c.Redirect(302, "/user/order?msg=ship_success")
}
