package front

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	BaseController
}

func (con ProductController) Category(c *gin.Context) {

	//测试获取动态路由的值
	id := c.Param("id")
	fmt.Println("----")
	fmt.Println(id)

	tpl := "itying/product/list.html"
	con.Render(c, tpl, gin.H{
		"page": 20,
	})

}
