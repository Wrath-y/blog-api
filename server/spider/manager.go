package spider

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Index(nextMarker string, page int) (oss.ListObjectsResult, error) {
	bucket, err := Bucket()
	marker := oss.Marker(nextMarker)
	list, err := bucket.ListObjects(oss.MaxKeys(page), marker)
	if err != nil {
		return list, err
	}

	return list, nil
}