package itying

import (
	"fmt"
	"goshopdemo/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	//1、获取顶部导航
	topNavList := []models.Nav{}
	if hasTopNavList := models.CacheDb.Get("topNavList", &topNavList); !hasTopNavList {
		models.DB.Where("status=1 AND position=1").Find(&topNavList)
		models.CacheDb.Set("topNavList", topNavList, 60*60)
		fmt.Println("数据库里面读取数据")
	} else {
		fmt.Println("redis里面读取数据")
	}
	//2、获取轮播图数据
	focusList := []models.Focus{}
	models.DB.Where("status=1 AND focus_type=1").Find(&focusList)
	//3、获取分类的数据
	goodsCateList := []models.GoodsCate{}
	//https://gorm.io/zh_CN/docs/preload.html
	models.DB.Where("pid = 0 AND status=1").Order("sort DESC").Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
		return db.Where("goods_cate.status=1").Order("goods_cate.sort DESC")
	}).Find(&goodsCateList)

	//4、获取中间导航
	middleNavList := []models.Nav{}
	models.DB.Where("status=1 AND position=2").Find(&middleNavList)

	for i := 0; i < len(middleNavList); i++ {
		relation := strings.ReplaceAll(middleNavList[i].Relation, "，", ",") //21，22,23,24
		relationIds := strings.Split(relation, ",")
		goodsList := []models.Goods{}
		models.DB.Where("id in ?", relationIds).Select("id,title,goods_img,price").Find(&goodsList)
		middleNavList[i].GoodsItems = goodsList
	}

	//方舟工业

	arkindustryList := models.GetGoodsByCategory(25, "best", 8)

	//配件

	otherList := models.GetGoodsByCategory(9, "all", 1)

	c.HTML(http.StatusOK, "itying/index/index.html", gin.H{
		"topNavList":      topNavList,
		"focusList":       focusList,
		"goodsCateList":   goodsCateList,
		"middleNavList":   middleNavList,
		"arkindustryList": arkindustryList,
		"otherList":       otherList,
	})

}
