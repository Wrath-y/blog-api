package spider

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

func Index(nextMarker string, page int) (oss.ListObjectsResult, error) {
	bucket, _ := Bucket()
	marker := oss.Marker(nextMarker)
	list, err := bucket.ListObjects(oss.MaxKeys(page), marker)
	if err != nil {
		return list, err
	}

	return list, nil
}

func Delete(name string) (int, error) {
	bucket, _ := Bucket()
	fmt.Println(name)
	err := bucket.DeleteObject(name)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func Count() (int, error) {
	bucket, _ := Bucket()
	var lsRes oss.ListObjectsResult
	var err error

	lsRes, err = bucket.ListObjects(oss.MaxKeys(1000))
	if err != nil {
		return 0, err
	}
	count := len(lsRes.Objects)

	return count, nil
}

func Bucket() (*oss.Bucket, error) {
	clientSer, _ := oss.New(viper.GetString("endPoint"),
		viper.GetString("accessKeyId"),
		viper.GetString("accessKeySecret"))

	// 获取存储空间。
	return clientSer.Bucket(viper.GetString("bucketName"))
}