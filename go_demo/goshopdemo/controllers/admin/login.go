package admin

import (
	"encoding/json"
	"fmt"
	"goshopdemo/models"
	"net/http"

	"github.com/gin-contrib/sessions"
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

	username := c.PostForm("username")
	password := c.PostForm("password")

	verifyValue := c.PostForm("verifyValue")

	//1.验证验证码
	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag {
		//2.查询数据库，判断用户密码是否存在
		userinfoList := []models.Manager{}
		password = models.Md5(password)
		models.DB.Where("username = ? and password = ?", username, password).Find(&userinfoList)

		if len(userinfoList) > 0 {

			//3.执行登录，保存用户信息，跳转到后台首页
			session := sessions.Default(c)
			//把结构体转换成josn
			userinSlice, _ := json.Marshal(userinfoList)
			session.Set("userinfo", string(userinSlice))
			session.Save()
			con.Success(c, "登录成功", "/admin")
		} else {
			con.Error(c, "用户名或密码错误", "/admin/login")
		}

	} else {
		con.Error(c, "验证码验证失败", "/admin/login")
	}

}

func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha(34, 100, 4)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con LoginController) LoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userinfo")
	session.Save()
	con.Success(c, "退出登录成功", "/admin/login")

}
