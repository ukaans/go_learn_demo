package models

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ginCookie struct{}

// 写入 Cookie
func (cookie ginCookie) Set(c *gin.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)

	desKey := []byte("itying.c") // 8位key
	encData, err := DesEncrypt(bytes, desKey)
	if err != nil {
		return
	}

	// 关键修复：使用 Base64 编码
	cookieValue := base64.RawURLEncoding.EncodeToString(encData)

	c.SetCookie(key, cookieValue, 3600*24*30, "/", "", false, true)
}

// 获取 Cookie
func (cookie ginCookie) Get(c *gin.Context, key string, obj interface{}) bool {
	valueStr, err1 := c.Cookie(key)
	if err1 != nil || valueStr == "" || valueStr == "[]" {
		return false
	}

	// Base64 解码
	raw, err := base64.RawURLEncoding.DecodeString(valueStr)
	if err != nil {
		c.SetCookie(key, "", -1, "/", "", false, true) // 清除坏Cookie
		return false
	}

	// DES 解密
	decData, err := DesDecrypt(raw, []byte("itying.c"))
	if err != nil {
		c.SetCookie(key, "", -1, "/", "", false, true)
		return false
	}

	err2 := json.Unmarshal(decData, obj)
	return err2 == nil
}

// 删除 Cookie
func (cookie ginCookie) Remove(c *gin.Context, key string) bool {
	c.SetCookie(key, "", -1, "/", "", false, true)
	return true
}

var Cookie = &ginCookie{}
