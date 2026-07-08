package front

import (
	"fmt"
	"goshopdemo/models"
	"net/http"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PassController struct {
	BaseController
}

// 获取验证码
func (con PassController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha(50, 120, 4)

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con PassController) Login(c *gin.Context) {
	//生成随机数
	fmt.Println(models.GetRandomNum())
	// c.HTML(http.StatusOK, "itying/pass/login.html", gin.H{})
	c.String(200, "login")
}
func (con PassController) RegisterStep1(c *gin.Context) {
	c.HTML(http.StatusOK, "itying/pass/register_step1.html", gin.H{})
}
func (con PassController) RegisterStep2(c *gin.Context) {
	c.HTML(http.StatusOK, "itying/pass/register_step2.html", gin.H{})
}
func (con PassController) RegisterStep3(c *gin.Context) {
	c.HTML(http.StatusOK, "itying/pass/register_step3.html", gin.H{})
}

func (con PassController) SendCode(c *gin.Context) {

	phone := c.Query("phone")
	verifyCode := c.Query("verifyCode")
	captchaId := c.Query("captchaId")
	fmt.Println(captchaId, verifyCode)
	// 1、验证图形验证码是否正确
	if flag := models.VerifyCaptcha(captchaId, verifyCode); !flag {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "验证码输入错误，请重试",
		})
		return
	}

	/*
		2、判断手机格式是否合法
				pattern := `^[\d]{11}$`
				reg := regexp.MustCompile(pattern)
				reg.MatchString(phone)
	*/
	pattern := `^[\d]{11}$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(phone) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号格式不合法",
		})
		return
	}

	//3、验证手机号是否注册过
	userList := []models.User{}
	models.DB.Where("phone = ?", phone).Find(&userList)
	if len(userList) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号已经注册，请直接登录",
		})
		return
	}
	//4、判断当前ip地址今天发送短信的次数

	ip := c.ClientIP()
	currentDay := models.GetDay() //20211211
	var sendCount int64
	models.DB.Table("user_temp").Where("ip=? AND add_day=?", ip, currentDay).Count(&sendCount)
	if sendCount > 4 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "此ip今天发送短信的次数已经达到上限，请明天再试",
		})
		return
	}
	//5、验证当前手机号今天发送的次数是否合法
	userTemp := []models.UserTemp{}
	smsCode := models.GetRandomNum()
	sign := models.Md5(phone + currentDay) //签名：主要用于页面跳转传值
	models.DB.Where("phone = ? AND add_day=?", phone, currentDay).Find(&userTemp)
	if len(userTemp) > 0 {
		if userTemp[0].SendCount > 2 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "此手机号今天发送短信的次数已经达到上限，请明天再试",
			})
			return
		} else {
			//1、生成短信验证码  发送验证码  调用前面课程的接口
			fmt.Println("----------自己集成发送短信的接口--------")
			fmt.Println(smsCode)
			//2、服务器保持验证码
			session := sessions.Default(c)
			session.Set("smsCode", smsCode)
			session.Save()

			//3、更新发送短信的次数
			oneUserTemp := models.UserTemp{}
			models.DB.Where("id=?", userTemp[0].Id).Find(&oneUserTemp)
			oneUserTemp.SendCount = oneUserTemp.SendCount + 1
			models.DB.Save(&oneUserTemp)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "发送短信成功",
				"sign":    sign,
			})
			return
		}

	} else {
		//1、生成短信验证码  发送验证码  调用前面课程的接口
		fmt.Println("----------自己集成发送短信的接口--------")
		fmt.Println(smsCode)
		//2、服务器保持验证码
		session := sessions.Default(c)
		session.Set("smsCode", smsCode)
		session.Save()

		//3、记录发送短信的次数

		oneUserTemp := models.UserTemp{
			Ip:        ip,
			Phone:     phone,
			SendCount: 1,
			AddDay:    currentDay,
			AddTime:   int(models.GetUnix()),
			Sign:      sign,
		}
		models.DB.Create(&oneUserTemp)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "发送短信成功",
			"sign":    sign,
		})
		return

	}

	c.JSON(200, gin.H{
		"SendCode": "SendCode",
	})

}
