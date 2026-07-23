package front

import (
	"goshopdemo/models"
	"math"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
}

func (con UserController) Index(c *gin.Context) {

	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	if user.Id == 0 {
		c.Redirect(302, "/")
		return
	}

	// 处理手机号脱敏
	phoneMask := "未绑定"
	if len(user.Phone) >= 11 {
		phoneMask = user.Phone[0:3] + "****" + user.Phone[7:11]
	}

	// 查询待支付订单数量
	var waitPayCount int64
	models.DB.Model(&models.Order{}).
		Where("uid = ? AND pay_status = 0", user.Id).
		Count(&waitPayCount)

	// 查询待收货订单数量
	var waitReceiveCount int64
	models.DB.Model(&models.Order{}).
		Where("uid = ? AND order_status = 3", user.Id).
		Count(&waitReceiveCount)

	con.Render(c, "itying/user/welcome.html", gin.H{
		"user":             user,
		"phoneMask":        phoneMask,
		"waitPayCount":     waitPayCount,
		"waitReceiveCount": waitReceiveCount,
	})
}
func (con UserController) OrderList(c *gin.Context) {
	// 当前页
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}
	//pageSize

	pageSize := 2

	//获取当前用户
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	//模糊查询
	where := "uid =" + models.String(user.Id)
	keywords := c.Query("keywords")
	if keywords != "" {
		//查询
		orderItemList := []models.OrderItem{}
		models.DB.Where("product_title like ?", "%"+keywords+"%").Find(&orderItemList)
		var str string
		// 字符串：   12,12,22
		for i := 0; i < len(orderItemList); i++ {
			if i == 0 {
				str += models.String(orderItemList[i].OrderId)
			} else {
				str += "," + models.String(orderItemList[i].OrderId)
			}
		}
		where += " AND id in (" + str + ")"
	}
	//按照状态筛选订单
	orderStatus, statusErr := models.Int(c.Query("orderStatus"))
	if statusErr == nil && orderStatus >= 0 {
		where += " AND order_status=" + models.String(orderStatus)
	} else {
		orderStatus = -1
	}

	//获取当前用户下面订单信息
	orderList := []models.Order{}
	models.DB.Where(where).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&orderList)

	//获取总数量
	var count int64
	models.DB.Where(where).Table("order").Count(&count)

	var tpl = "itying/user/order.html"
	con.Render(c, tpl, gin.H{
		"order":       orderList,
		"page":        page,
		"keywords":    keywords,
		"orderStatus": orderStatus,
		"totalPages":  math.Ceil(float64(count) / float64(pageSize)),
	})
}

func (con UserController) OrderInfo(c *gin.Context) {

	id, err := models.Int(c.Query("id"))
	if err != nil {
		c.Redirect(302, "/user/order")
	}
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	order := []models.Order{}
	models.DB.Where("id=? And uid=?", id, user.Id).Preload("OrderItem").Find(&order)

	if len(order) == 0 {
		c.Redirect(302, "/user/order")
		return
	}
	var tpl = "itying/user/order_info.html"
	con.Render(c, tpl, gin.H{
		"order": order[0],
	})
}
