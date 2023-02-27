package pixivoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

func List(nextMarker string, page int) (*oss.ListObjectsResult, error) {
	bucket, err := Bucket()
	if err != nil {
		return nil, err
	}
	list, err := bucket.ListObjects(oss.MaxKeys(page), oss.Marker(nextMarker))
	if err != nil {
		return &list, err
	}

	return &list, nil
}

func Delete(name string) (int, error) {
	bucket, _ := Bucket()
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
	clientSer, _ := oss.New(viper.GetString("aliyun.oss.endpoint"),
		viper.GetString("aliyun.oss.access_keyid"),
		viper.GetString("aliyun.oss.access_keysecret"))
	// 获取存储空间。
	return clientSer.Bucket(viper.GetString("aliyun.oss.bucket_pixiv"))
}
