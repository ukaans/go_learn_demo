package models

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

//配置store
// var store = base64Captcha.DefaultMemStore

// 配置redisStore
var store base64Captcha.Store = RedisStore{}

// 获取验证码
func MakeCaptcha() (string, string, error) {
	var driver base64Captcha.Driver
	driverString := base64Captcha.DriverString{
		Height:          40,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}
	//ConvertFonts按名称加载字体
	driver = driverString.ConvertFonts()

	c := base64Captcha.NewCaptcha(driver, store)
	//Generate生成随机id、base64图像字符串
	id, b64s, _, err := c.Generate()
	return id, b64s, err
}

// 验证验证码
func VerifyCaptcha(id string, capt string) bool {
	if store.Verify(id, capt, true) {
		return true
	} else {
		return false
	}
}
