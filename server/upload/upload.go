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

func GetSign(c *gin.Context) interface{} {
	conditions := make([]map[int]string, 2)
	conditions[0] = make(map[int]string)
	conditions[1] = make(map[int]string)
	conditions[0][0] = "content-length-range"
	conditions[0][1] = "0"
	conditions[0][2] = "1048576000"
	conditions[1][0] = "starts-with"
	conditions[1][1] = "$key"
	conditions[1][2] = ""


	now := time.Now()
	expire, _ := time.ParseDuration("30s")
	after := now.Add(expire)

	arr := make(map[string]interface{})
	arr["expiration"] = after.Format("2006-01-02T15:04:05Z")
	arr["conditions"] = conditions

	return arr
}

func Create(c *gin.Context) {
	panic(c.Request.Body)
	objectName := c.Param("title") + time.Now().Format("2006-01-02 15:04:05")
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