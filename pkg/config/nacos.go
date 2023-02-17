package config

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type api struct {
	loginApi  string
	getApi    string
	listenApi string
}

const (
	contentType      = "application/x-www-form-urlencoded;charset=utf-8"
	splitConfig      = string(rune(1))
	splitConfigInner = string(rune(2))
)

type logger interface {
	Info(message string, request, response interface{}, t ...time.Time)
	ErrorL(message string, request, response interface{}, t ...time.Time)
	SetV1(v1 string)
	SetV2(v2 string)
	SetV3(v3 string)
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type NacosConfig struct {
	HttpClient  httpClient
	ServerAddr  string
	accessToken string
	tokenTTL    int
	Username    string
	Password    string
	Logger      logger
	PollTime    time.Duration
	NeedInit    chan struct{}
	api         *api
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenTTL    int    `json:"tokenTtl"`
	GlobalAdmin bool   `json:"globalAdmin"`
}

func NewNacosConfig(options ...func(c *NacosConfig)) *NacosConfig {
	nc := &NacosConfig{
		HttpClient: http.DefaultClient,
		PollTime:   10 * time.Second,
	}

	for _, option := range options {
		option(nc)
	}

	if nc.Username != "" && nc.Password != "" {
		nc.api = &api{
			loginApi:  "/nacos/v1/auth/login",
			getApi:    "/nacos/v1/cs/configs",
			listenApi: "/nacos/v1/cs/configs/listener",
		}
		if err := nc.login(); err != nil {
			panic(err)
		}
	}

	return nc
}

func (n *NacosConfig) login() error {
	n.Logger.SetV3("login")
	////n.Logger.Info(fmt.Sprintf("nacos login server:[%s:%s]", n.ServerAddr, n.Username), nil, nil)

	v := url.Values{}
	v.Add("username", n.Username)
	v.Add("password", n.Password)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", n.ServerAddr, n.api.loginApi), strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentType)

	resp, err := n.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bb, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("nacos login fail:%s", string(bb))
	}

	loginResp := &LoginResponse{}

	if err := json.Unmarshal(bb, loginResp); err != nil {
		return err
	}
	n.accessToken = loginResp.AccessToken
	n.tokenTTL = loginResp.TokenTTL - 600

	return nil
}

func (n *NacosConfig) Get(namespace, group, dataId string) (string, error) {
	n.Logger.SetV3("Get")
	//n.Logger.Info(fmt.Sprintf("nacos get config:[namespace:%s,group:%s,dataId:%s]", namespace, group, dataId), nil, nil)

	v := url.Values{}
	v.Add("tenant", namespace)
	v.Add("group", group)
	v.Add("dataId", dataId)
	if n.accessToken != "" {
		v.Add("accessToken", n.accessToken)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s?", n.ServerAddr, n.api.getApi)+v.Encode(), nil)
	if err != nil {
		return "", err
	}

	timeStamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("timeStamp", timeStamp)

	resp, err := n.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bb, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("nacos get fail:%s", string(bb))
	}

	return string(bb), nil
}

func (n *NacosConfig) ListenAsync(namespace, group, dataId string, fn func(cnf string)) {
	n.Logger.SetV3("ListenAsync")
	ret, err := n.Get(namespace, group, dataId)
	if err != nil {
		panic(err)
	}

	contentMd5 := md5string(ret)

	go func() {
		t1 := time.NewTicker(time.Duration(n.tokenTTL) * time.Second)
		t2 := time.NewTicker(n.PollTime)
		for {
			select {
			// token到期刷新
			case <-t1.C:
				if n.Username != "" {
					if err := n.login(); err != nil {
						n.Logger.ErrorL("登录失败", nil, err)
					}
				}
			// 每10秒监听配置
			case <-t2.C:
				update, err := n.Listen(namespace, group, dataId, contentMd5)
				if err != nil {
					n.Logger.ErrorL("监听配置失败", nil, err)
					continue
				}
				if update {
					ret, err := n.Get(namespace, group, dataId)
					if err != nil {
						n.Logger.ErrorL("获取配置失败", nil, err)
						continue
					}

					contentMd5 = md5string(ret)
					//n.Logger.Info(fmt.Sprintf("nacos listen refresh:[namespace:%s,group:%s,dataId:%s,md5:%s]", namespace, group, dataId, contentMd5), nil, nil)
					fn(ret)
				}
			case <-n.NeedInit:
				ret, err := n.Get(namespace, group, dataId)
				if err != nil {
					n.Logger.ErrorL("获取配置失败", nil, err)
					continue
				}

				contentMd5 = md5string(ret)
				//n.Logger.Info(fmt.Sprintf("nacos listen refresh:[namespace:%s,group:%s,dataId:%s,md5:%s]", namespace, group, dataId, contentMd5), nil, nil)
				fn(ret)
			}
		}
	}()
}

func (n *NacosConfig) Listen(namespace, group, dataId, md5 string) (bool, error) {
	//n.Logger.Info(fmt.Sprintf("nacos listen start:[namespace:%s,group:%s,dataId:%s,md5:%s]", namespace, group, dataId, md5), nil, nil)

	content := dataId + splitConfigInner + group + splitConfigInner + md5 + splitConfigInner + namespace + splitConfig

	v := url.Values{}
	if n.Username != "" {
		v.Add("Listening-Configs", content)
	} else {
		v.Add("Probe-Modify-Request", content)
	}
	v.Add("tenant", namespace)
	if n.accessToken != "" {
		v.Add("accessToken", n.accessToken)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", n.ServerAddr, n.api.listenApi), strings.NewReader(v.Encode()))
	if err != nil {
		return false, err
	}

	timeStamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	req.Header.Add("Long-Pulling-Timeout", "3000")
	req.Header.Add("User-Agent", "Nacos-go-client/v1.0.1")
	req.Header.Add("exConfigInfo", "true")
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("timeStamp", timeStamp)

	resp, err := n.HttpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	bb, err := io.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("nacos listen response error:%s", string(bb))
	}

	str := strings.Split(string(bb), "%02")

	// 如果返回数据不为空则代表有变化的文件
	if len(str) > 0 && str[0] == dataId {
		return true, nil
	}

	return false, nil
}

func md5string(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func signSha1(encryptText, encryptKey string) string {
	// hmac ,use sha1
	key := []byte(encryptKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(encryptText))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
