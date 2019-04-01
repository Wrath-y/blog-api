package uploadController

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-blog/server/errno"
	"go-blog/struct"
	"hash"
	"time"
)

type ConfigStruct struct{
	Expiration string `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct{
	AccessKeyId string `json:"accessid"`
	Host string `json:"host"`
	Expire int64 `json:"expire"`
	Signature string `json:"signature"`
	Policy string `json:"policy"`
	Directory string `json:"dir"`
	Callback string `json:"callback"`
}

type Data struct{
	AccessUrl string `json:"access_url"`
	Drive string `json:"drive"`
	FileField string `json:"file_field"`
	Form interface{} `json:"form"`
	Headers interface{} `json:"headers"`
	UploadUrl string `json:"upload_url"`
	Policy string `json:"policy"`
	OSSAccessKeyId string `json:"oss_access_key_id"`
	Signature string `json:"signature"`
	SuccessActionStatus int `json:"success_action_status"`
}

type Form struct {
	OSSAccessKeyId string `json:"OSSAccessKeyId"`
	Policy string `json:"policy"`
	Signature string `json:"signature"`
	SuccessActionStatus int `json:"success_action_status"`
}

type CallbackParam struct{
	CallbackUrl string `json:"callbackUrl"`
	CallbackBody string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func GetGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func Index(c *gin.Context) {
	accessKeyId := viper.GetString("accessKeyId")
	accessKeySecret :=viper.GetString("accessKeySecret")
	// host的格式为 bucketname.endpoint
	host := "http://blog-ico." + viper.GetString("endPoint")
	// 上传文件时指定的前缀。
	uploadDir := ""
	expireTime := 30
	now := time.Now().Unix()
	expireEnd := now + int64(expireTime)
	var tokenExpire = GetGmtIso8601(expireEnd)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition= append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	h.Write([]byte(debyte))
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expireEnd
	policyToken.Signature = string(signedStr)
	policyToken.Directory = uploadDir
	policyToken.Policy = string(debyte)

	//policy, err := json.Marshal(policyToken)
	if err != nil {
		_struct.Response(c, errno.UploadError, err)

		return
	}

	response := &Data{
		AccessUrl: policyToken.Host,
		Drive: "oss",
		FileField: "file",
		OSSAccessKeyId: accessKeyId,
		Policy: policyToken.Policy,
		Signature: policyToken.Signature,
		SuccessActionStatus: 200,
		Headers: []int{},
		UploadUrl: policyToken.Host,
	}
	_struct.Response(c, nil, response)

	return
}