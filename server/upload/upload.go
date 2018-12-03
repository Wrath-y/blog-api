package upload

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-blog/server/errno"
	"go-blog/struct"
	"time"
)

func Version() string {
	return "OSS Go SDK Version: " + oss.Version
}

func Create(c *gin.Context) {
	panic(c.Request.Body)
	objectName := c.Param("title") + time.Now().Format("2006-01-02_15:04:05")
	localFileName := "<yourLocalFileName>"
	// 创建OSSClient实例。
	client, err := oss.New(viper.GetString("endpoint"), viper.GetString("accessKeyId"), viper.GetString("accessKeySecret"))
	if err != nil {
		_struct.Response(c, errno.UploadError, nil)

		return
	}
	// 获取存储空间。
	bucket, err := client.Bucket(viper.GetString("bucketName"))
	if err != nil {
		_struct.Response(c, errno.UploadError, nil)

		return
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		_struct.Response(c, errno.UploadError, nil)

		return
	}
}