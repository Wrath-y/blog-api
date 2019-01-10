package spider

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-blog/server/errno"
	"go-blog/struct"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Img struct {
	ImgId  		string
	Title  		string
	Url	    	string
	Praise 		string
}

var waitGroup = sync.WaitGroup{}
var errChan chan error
var errCodeChan chan error
var finished chan bool
var lock = new(sync.Mutex)

func Get(c *gin.Context) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: tr,
		Jar: jar,
		Timeout: time.Second * 60,
	}
	var loginResp *http.Response
	loginReq, err := http.NewRequest("GET", loginRrl, nil)
	loginResp, err = client.Do(loginReq)
	if err != nil || loginResp == nil {
		_struct.Response(c, errno.ErrCurl.Add("login"), err)
		return
	}
	defer loginResp.Body.Close()
	loginBody, err := ioutil.ReadAll(loginResp.Body)

	exp := "post_key\" value=\"(.+?)\">"
	r, _ := regexp.Compile(exp)
	postKey := r.FindStringSubmatch(string(loginBody))[1]

	value := url.Values{}
	value.Add("pixiv_id", pixivId)
	value.Add("password", password)
	value.Add("post_key", postKey)
	value.Add("source", "pc")
	value.Add("ref", ref)
	value.Add("return_to", returnTo)
	form := ioutil.NopCloser(strings.NewReader(value.Encode()))

	postLoginReq, _ := http.NewRequest("POST", loginPostRrl, form)
	postLoginReq.Header.Set("Content-Type","application/x-www-form-urlencoded")
	postLoginReq.Header.Set("Accept",accept)
	postLoginReq.Header.Set("Accept-Encoding", "deflate, br")
	postLoginReq.Header.Set("Origin", origin)
	postLoginReq.Header.Set("Referer", referer)
	postLoginReq.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(postLoginReq)
	if err != nil || resp == nil {
		_struct.Response(c, errno.ErrCurl.Add("loginpost"), err)
		return
	}
	GetList(c, client)
	return
}

func GetList(c *gin.Context, client *http.Client)  {
	bookmarkReq, _ := http.NewRequest("GET", bookmark, nil)
	bookmarkResp, err := client.Do(bookmarkReq)
	if err != nil || bookmarkResp == nil {
		_struct.Response(c, errno.ErrCurl.Add("bookmark"), err)
		return
	}
	var buf []byte
	buf, _ = ioutil.ReadAll(bookmarkResp.Body)
	content := string(buf)
	allContent := content

	pageExpInfos, _ := regexp.Compile(`w&amp;p=\d+[\s\S]*(\d+)[\s\S]*s="next"`)
	page, _ := strconv.Atoi(pageExpInfos.FindStringSubmatch(content)[1])
	if page == 0 {
		page = 1
	}
	p := 1
	for {
		if p > 1 {
			bookmarkReq, _ = http.NewRequest("GET", bookmark + "?rest=show&p=" + strconv.Itoa(p), nil)
			bookmarkResp, err = client.Do(bookmarkReq)
			if bookmarkResp == nil {
				_struct.Response(c, errno.ErrCurl.Add("bookmarkwithpage"), err)
				return
			}
			buf, _ = ioutil.ReadAll(bookmarkResp.Body)
			content = string(buf)
			allContent += content
			pageExpInfos, _ = regexp.Compile(`w&amp;p=\d+[\s\S]*s="">(.+?)<[\s\S]*s="next"`)
			page, _ = strconv.Atoi(pageExpInfos.FindStringSubmatch(content)[1])
			if page == 0 {
				page = 1
			}
		}
		p = p + 1
		if p > page {
			break
		}
	}
	defer bookmarkResp.Body.Close()
	size := (page + 1) * 20
	k := 0
	imgSlice := make([]Img, size)
	r, _ := regexp.Compile(`data-id="(.+?)".+?title="(.+?)".+?e"></i>(.+?)</a>`)
	imgExpInfos := r.FindAllStringSubmatch(allContent, size)

	errChan = make(chan error)
	errCodeChan = make(chan error)
	finished = make(chan bool, 1)
	for _, v := range imgExpInfos {
		imgSlice[k].ImgId = v[1]
		imgSlice[k].Title = v[2]
		imgSlice[k].Url = "https://www.pixiv.net/member_illust.php?mode=medium&illust_id=" + v[1]
		imgSlice[k].Praise = v[3]
		waitGroup.Add(1)
		go GetDetail(c, client, imgSlice[k], false)
		k++
	}
	go func() {
		waitGroup.Wait()
		close(finished)
		close(errChan)
		close(errCodeChan)
	}()
	select {
		case <- finished:
		case err := <- errChan:
			errCode := <- errCodeChan
			_struct.Response(c, errCode, err)
			return
	}

	return
}

func GetDetail(c *gin.Context, client *http.Client, img Img, try bool) {
	defer waitGroup.Done()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	lock.Lock()
	defer lock.Unlock()
	req, _ := http.NewRequest("GET", img.Url, nil)
	res, err := client.Do(req)
	if err != nil || res == nil {
		errChan <- err
		errCodeChan <- errno.ErrCurl.Add("imgurl")
		return
	}
	defer res.Body.Close()

	var buf []byte
	var content string
	buf, _ = ioutil.ReadAll(res.Body)
	exp, err := regexp.Compile(`nal":"(.+?)"}`)
	contentArr := exp.FindStringSubmatch(string(buf))
	if len(contentArr) > 1 {
		content = contentArr[1]
	} else {
		errChan <- errno.ErrExp.Add("未匹配到网页内容")
		errCodeChan <- errno.ErrExp
		return
	}
	exp, _ = regexp.Compile(`\\`)
	src := exp.ReplaceAllString(content, "")
	var suffix string
	exp, err = regexp.Compile(`p0(.+)`)
	suffixArr := exp.FindStringSubmatch(src)
	if len(suffixArr) > 1 {
		suffix = suffixArr[1]
	} else {
		errChan <- errno.ErrExp.Add("未匹配到类型后缀")
		errCodeChan <- errno.ErrExp
		return
	}

	exp, _ = regexp.Compile(`/`)
	effecTitle := exp.ReplaceAllString(img.Title, "-")

	img.Title = effecTitle

	if Exist("static/pixiv/" + img.Title + suffix) {
		return
	}

	imgreq, _ := http.NewRequest("GET", src, nil)
	imgreq.Header.Set("Accept",accept)
	imgreq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	imgreq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7")
	imgreq.Header.Set("Referer", img.Url)
	imgreq.Header.Set("pragma", "no-cache")
	imgreq.Header.Set("Cache-Control", "no-cache")
	imgreq.Header.Set("User-Agent", userAgent)
	if try {
		var createDate string
		exp, _ = regexp.Compile(`createDate":"(.+?)",`)
		createDateArr := exp.FindStringSubmatch(string(buf))
		if len(createDateArr) > 1 {
			createDate = createDateArr[1]
		} else {
			errChan <- errno.ErrExp.Add("未匹配到创建时间")
			errCodeChan <- errno.ErrExp
			return
		}
		exp, _ = regexp.Compile(`T`)
		effCreateDate := exp.ReplaceAllString(createDate, " ")
		exp, _ = regexp.Compile(`\+.+`)
		date := exp.ReplaceAllString(effCreateDate, "")
		timestamp, _ := time.Parse("2006-01-02 15:04:05", date)
		GMTtime := timestamp.Format("Mon, 02 Jan 2006 15:04:05 GMT")
		imgreq.Header.Set("Upgrade-Insecure-Requests", "1")
		imgreq.Header.Set("If-Modified-Since", GMTtime)
	}

	imgRes, err := client.Do(imgreq)
	if imgRes.ContentLength > 0 {
		if err != nil || imgRes == nil {
			errChan <- err
			errCodeChan <- errno.ErrCurl.Add("imgres")
			return
		}
		defer imgRes.Body.Close()

		imgBytes, err := ioutil.ReadAll(imgRes.Body)
		if err != nil {
			errChan <- err
			errCodeChan <- errno.ErrIoutilReadAll
			return
		}

		// 获取存储空间。
		bucket, err := Bucket()
		if err != nil {
			errChan <- err
			errCodeChan <- errno.UploadError.Add("获取储存空间失败")
			return
		}
		isExist, err := bucket.IsObjectExist(img.Title + suffix)
		if err != nil {
			errChan <- err
			errCodeChan <- errno.UploadError.Add("判断图片是否存在失败")
			return
		}
		fmt.Println(img.Title, isExist)
		if isExist == true {
			fmt.Println(img.Title + "已存在")
			return
		}
		err = bucket.PutObject(img.Title + suffix, bytes.NewReader(imgBytes))
		if err != nil {
			errChan <- err
			errCodeChan <- errno.UploadError.Add("上传byte数组失败")
			return
		}
		fmt.Println(img.Title + "上传成功")

		newFile, err := os.Create("static/pixiv/" + img.Title + suffix)
		if err != nil {
			errChan <- err
			errCodeChan <- errno.ErrOsCreate
			return
		}
		defer newFile.Close()
		w, err := io.Copy(newFile, bytes.NewReader(imgBytes))
		if w == 0 || err != nil  {
			errChan <- err
			errCodeChan <- errno.ErrIoCopy
			return
		}
	} else {
		GetDetail(c, client, img, true)
	}
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Bucket() (*oss.Bucket, error) {
	clientSer, _ := oss.New(viper.GetString("endPoint"),
		viper.GetString("accessKeyId"),
		viper.GetString("accessKeySecret"))

	// 获取存储空间。
	return clientSer.Bucket(viper.GetString("bucketName"))
}