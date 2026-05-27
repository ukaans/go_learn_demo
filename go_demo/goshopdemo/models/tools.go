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
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	. "github.com/hunterhug/go_image"
)

// ... 保留原有的时间、加密等工具函数 ...

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

// getOssClient 从环境变量读取 OSS 凭证并初始化客户端
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
		bucketName = "golearndemo" // 默认 bucket 名称
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

// GetSettingFromColumn 从数据库 setting 表获取指定字段的值
func GetSettingFromColumn(columnName string) string {
	setting := Setting{}
	DB.First(&setting)
	v := reflect.ValueOf(setting)
	val := v.FieldByName(columnName).String()
	return val
}

// GetOssStatus 从环境变量读取 OSS 开关状态（1 开启，0 关闭）
func GetOssStatus() int {
	statusStr := os.Getenv("OSS_STATUS")
	if statusStr == "" {
		return 0 // 默认关闭
	}
	status, err := Int(statusStr)
	if err != nil {
		return 0
	}
	return status
}

func FormatImg(str string) string {
	ossStatus := GetOssStatus()
	if ossStatus == 1 {
		return GetSettingFromColumn("") + str
	}
	return "/" + str
}

// UploadImg 根据 OSS 状态选择上传方式
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

// ResizeGoodsImage 生成商品缩略图（保留原逻辑）
func ResizeGoodsImage(filename string) {
	extname := path.Ext(filename)
	ThumbnailSize := strings.ReplaceAll(GetSettingFromColumn("ThumbnailSize"), "，", ",")
	thumbnailSizeSlice := strings.Split(ThumbnailSize, ",")
	for i := 0; i < len(thumbnailSizeSlice); i++ {
		savepath := filename + "_" + thumbnailSizeSlice[i] + "x" + thumbnailSizeSlice[i] + extname
		w, _ := Int(thumbnailSizeSlice[i])
		err := ThumbnailF2F(filename, savepath, w, w)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Sub(a int, b int) int {
	return a - b
}
