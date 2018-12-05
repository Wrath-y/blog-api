package upload

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"hash"
	"strconv"
	"time"
)

// 请填写您的AccessKeyId。
var accessKeyId string = viper.GetString("accessKeyId")
// 请填写您的AccessKeySecret。
var accessKeySecret string = viper.GetString("accessKeySecret")
// host的格式为 bucketname.endpoint ，请替换为您的真实信息。
var host string = viper.GetString("host")
// callbackUrl为 上传回调服务器的URL，请将下面的IP和Port配置为您自己的真实信息。
var callbackUrl string = "http://88.88.88.88:8888";
// 用户上传文件时指定的前缀。
var upload_dir string = ""
var expire_time int64 = 30

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

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

type CallbackParam struct{
	CallbackUrl string `json:"callbackUrl"`
	CallbackBody string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func GetPolicyToken() string {
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = get_gmt_iso8601(expire_end)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))

	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		fmt.Println("callback json err:", err)
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	var policyToken PolicyToken
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = upload_dir
	policyToken.Policy = string(debyte)
	policyToken.Callback = string(callbackBase64)
	response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
	}

	return string(response)
}

type EscapeError string
func (e EscapeError) Error() string {
	return "invalid URL escape " + strconv.Quote(string(e))
}

type InvalidHostError string
func (e InvalidHostError) Error() string {
	return "invalid character " + strconv.Quote(string(e)) + " in host name"
}
