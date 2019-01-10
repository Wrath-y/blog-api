package spider

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Index(page int) ([]oss.ListObjectsResult, error) {
	var list []oss.ListObjectsResult
	bucket, err := Bucket()
	if err != nil {
		return list, err
	}
	marker := oss.Marker("")
	for {
		lsRes, err := bucket.ListObjects(oss.MaxKeys(page), marker)
		if err != nil {
			return list, err
		}
		list = append(list, lsRes)
		marker = oss.Marker(lsRes.NextMarker)

		fmt.Println("Objects:", lsRes.Objects)

		if !lsRes.IsTruncated {
			break
		}
	}

	return list, nil
}