package admin

import (
	"fmt"
	"goshopdemo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login/login.html", gin.H{})

}
func (con LoginController) DoLogin(c *gin.Context) {

	captchaId := c.PostForm("captchaId")

	verifyValue := c.PostForm("verifyValue")

	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag == true {
		c.String(http.StatusOK, "验证码成功")
	} else {
		c.String(http.StatusOK, "验证码失败")
	}

}

func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con LoginController) VerifyCaptcha(c *gin.Context) {
	c.String(http.StatusOK, "-add--文章-")
}
