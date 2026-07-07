package routers

import (
	"goshopdemo/controllers/front"

	"github.com/gin-gonic/gin"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", front.DefaultController{}.Index)

		defaultRouters.GET("/category:id", front.ProductController{}.Category)
		defaultRouters.GET("/detail", front.ProductController{}.Detail)
		defaultRouters.GET("/product/getImgList", front.ProductController{}.GetImgList)

		defaultRouters.GET("/cart", front.CartController{}.Get)
		defaultRouters.GET("/cart/addCart", front.CartController{}.AddCart)

		defaultRouters.GET("/cart/successTip", front.CartController{}.AddCartSuccess)
		defaultRouters.GET("/cart/decCart", front.CartController{}.DecCart)
		defaultRouters.GET("/cart/incCart", front.CartController{}.IncCart)

		defaultRouters.GET("/cart/changeOneCart", front.CartController{}.ChangeOneCart)
		defaultRouters.GET("/cart/changeAllCart", front.CartController{}.ChangeAllCart)
		defaultRouters.GET("/cart/delCart", front.CartController{}.DelCart)

	}
}
