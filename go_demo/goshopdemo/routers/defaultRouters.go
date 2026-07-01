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

	}
}
