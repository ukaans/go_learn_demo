package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

// ========== 时间工具函数保持不变 ==========

func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func DateToUnix(str string) int64 {
	layout := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func GetUnix() int64     { return time.Now().Unix() }
func GetUnixNano() int64 { return time.Now().UnixNano() }

func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetDay() string {
	return time.Now().Format("20060102")
}

func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Int(str string) (int, error)       { return strconv.Atoi(str) }
func Float(str string) (float64, error) { return strconv.ParseFloat(str, 64) }
func String(n int) string               { return strconv.Itoa(n) }
func Str2Html(str string) template.HTML { return template.HTML(str) }

var (
	ossClient     *oss.Client
	ossClientOnce sync.Once
	ossInitErr    error
)

// 凭证从环境变量读取
func getOssClient() (*oss.Client, error) {
	ossClientOnce.Do(func() {
		endpoint := os.Getenv("OSS_ENDPOINT")
		accessKeyID := os.Getenv("OSS_ACCESS_KEY_ID")
		accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")

		if endpoint == "" || accessKeyID == "" || accessKeySecret == "" {
			ossInitErr = errors.New("OSS configuration missing: please set OSS_ENDPOINT, OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET environment variables")
			return
		}

		ossClient, ossInitErr = oss.New(endpoint, accessKeyID, accessKeySecret)
	})
	return ossClient, ossInitErr
}

// OssUpload 上传文件到 OSS
func OssUpload(file *multipart.FileHeader, dst string) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open upload file: %w", err)
	}
	defer f.Close()

	client, err := getOssClient()
	if err != nil {
		return "", fmt.Errorf("OSS client init failed: %w", err)
	}

	bucketName := os.Getenv("OSS_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "golearndemo"
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to get OSS bucket: %w", err)
	}

	if err := bucket.PutObject(dst, f); err != nil {
		return "", fmt.Errorf("failed to upload to OSS: %w", err)
	}

	return dst, nil
}

func GetSettingFromColumn(columnName string) string {
	setting := Setting{}
	DB.First(&setting)
	v := reflect.ValueOf(setting)
	val := v.FieldByName(columnName).String()
	return val
}

func GetOssStatus() int {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}
	ossStatus, _ := Int(config.Section("oss").Key("status").String())
	return ossStatus
}

func FormatImg(str string) string {
	ossStatus := GetOssStatus()
	if ossStatus == 1 {
		return GetSettingFromColumn("") + str
	}
	return "/" + str
}

func UploadImg(c *gin.Context, picName string) (string, error) {
	ossStatus := GetOssStatus()
	if ossStatus == 1 {
		return OssUploadImg(c, picName)
	}
	return LocalUploadImg(c, picName)
}

func OssUploadImg(c *gin.Context, picName string) (string, error) {
	file, err := c.FormFile(picName)
	if err != nil {
		return "", err
	}

	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg": true, ".png": true, ".gif": true, ".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}
	day := GetDay()
	dir := "static/upload/" + day
	fileName := strconv.FormatInt(GetUnixNano(), 10) + extName
	dst := path.Join(dir, fileName)

	if _, err := OssUpload(file, dst); err != nil {
		return "", err
	}
	return dst, nil
}

func LocalUploadImg(c *gin.Context, picName string) (string, error) {
	file, err := c.FormFile(picName)
	if err != nil {
		return "", err
	}

	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg": true, ".png": true, ".gif": true, ".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}

	day := GetDay()
	dir := "./static/upload/" + day
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	fileName := strconv.FormatInt(GetUnixNano(), 10) + extName
	dst := path.Join(dir, fileName)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	return dst, nil
}
