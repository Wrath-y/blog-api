package config

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

var ErrInvalidConf = errors.New("nacos缺少必要配置")

const (
	DefaultPollTime = 30 * time.Second
	DefaultTimeout  = 15 * time.Second
)

type nacosParams struct {
	address     string
	username    string
	password    string
	dataID      string
	group       string
	namespaceID string
	pollTime    time.Duration
	timeout     time.Duration
}

func loadNacosParams() (*nacosParams, error) {
	conf := &nacosParams{
		address:     viper.GetString("nacos.address"),
		username:    viper.GetString("nacos.username"),
		password:    viper.GetString("nacos.password"),
		dataID:      viper.GetString("nacos.data_id"),
		group:       viper.GetString("nacos.group"),
		namespaceID: viper.GetString("nacos.namespace_id"),
		pollTime:    viper.GetDuration("nacos.poll_time"),
		timeout:     viper.GetDuration("nacos.timeout"),
	}

	if conf.address == "" || conf.username == "" || conf.password == "" || conf.dataID == "" || conf.group == "" || conf.namespaceID == "" {
		return nil, ErrInvalidConf
	}
	if conf.pollTime == 0 {
		conf.pollTime = DefaultPollTime
	}
	if conf.timeout == 0 {
		conf.timeout = DefaultTimeout
	}

	return conf, nil
}

func ListenNacos(l logger, httpClient httpClient, callbacks ...func(cnf string)) {
	l.SetV1("config")
	l.SetV2("nacos_sync")
	l.SetV3("ListenNacos")
	nacosParams, err := loadNacosParams()
	if err != nil {
		l.ErrorL("[nacos] 加载配置失败: %s", nil, err.Error())
		return
	}

	nacosConf := NewNacosConfig(func(c *NacosConfig) {
		c.ServerAddr = nacosParams.address
		c.Username = nacosParams.username
		c.Password = nacosParams.password
		c.PollTime = nacosParams.pollTime
		c.HttpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: nacosParams.timeout,
		}

		c.Logger = l
		c.HttpClient = httpClient
	})
	nacosConf.ListenAsync(nacosParams.namespaceID, nacosParams.group, nacosParams.dataID, func(cnf string) {
		l.Info("[nacos] 监听到配置文件有改变，开始获取", nacosParams, nil)

		content, err := nacosConf.Get(nacosParams.namespaceID, nacosParams.group, nacosParams.dataID)
		if err != nil {
			l.ErrorL("[nacos] 获取最新配置失败: %s", nacosParams, err.Error())
			return
		}
		if content == "" {
			l.ErrorL("[nacos] 获取到最新的配置为空", nacosParams, nil)
			return
		}

		// 同步到本地配置文件，之后会被viper监听到并重新加载
		if err := writeFile(DefaultRelationPath, content); err != nil {
			l.ErrorL("[nacos] 更新配置文件失败: %s", nacosParams, err.Error())
		}
		l.Info("[nacos] 更新配置文件成功\n%s", nacosParams, content)

		// 执行callback
		for _, callbackFunc := range callbacks {
			callbackFunc(cnf)
		}
	})
}

func writeFile(configPath, configContent string) (err error) {
	// 打开配置文件
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "配置文件打开失败")
	}
	defer file.Close()

	// 阻塞模式下，加排他锁
	//if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
	//	return errors.Wrapf(err, "文件加锁失败")
	//}
	//defer func() {
	//	if err = syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
	//		err = errors.Wrapf(err, "文件解锁失败")
	//	}
	//}()

	// 加载配置信息
	_, err = file.WriteString(configContent)
	if err != nil {
		return errors.Wrapf(err, "写入配置文件失败")
	}

	return nil
}
