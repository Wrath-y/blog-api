package spider

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
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

var waitGroup = new(sync.WaitGroup)

func Login(c *gin.Context) {
	proxy, _ := url.Parse("http://127.0.0.1:8123")
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
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
	}
	GetList(c, client)
	return
}

func GetList(c *gin.Context, client *http.Client)  {
	bookmarkReq, _ := http.NewRequest("GET", bookmark, nil)
	bookmarkResp, err := client.Do(bookmarkReq)
	if err != nil || bookmarkResp == nil {
		_struct.Response(c, errno.ErrCurl.Add("bookmark"), err)
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
		fmt.Println(p)
		if p > 1 {
			bookmarkReq, _ = http.NewRequest("GET", bookmark + "?rest=show&p=" + strconv.Itoa(p), nil)
			bookmarkResp, err = client.Do(bookmarkReq)
			if bookmarkResp == nil {
				_struct.Response(c, errno.ErrCurl.Add("bookmarkwithpage"), err)
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

	for _, v := range imgExpInfos {
		imgSlice[k].ImgId = v[1]
		imgSlice[k].Title = v[2]
		imgSlice[k].Url = "https://www.pixiv.net/member_illust.php?mode=medium&illust_id=" + v[1]
		imgSlice[k].Praise = v[3]
		waitGroup.Add(1)
		GetDetail(c, client, imgSlice[k])
		k++
		if k > 3 {
			break
		}
	}
	// waitGroup.Wait()
	return
}

func GetDetail(c *gin.Context, client *http.Client, img Img) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	req, _ := http.NewRequest("GET", img.Url, nil)
	res, err := client.Do(req)
	if err != nil || res == nil {
		_struct.Response(c, errno.ErrCurl.Add("imgurl"), err)
	}
	defer res.Body.Close()

	var buf []byte
	var content string
	buf, _ = ioutil.ReadAll(res.Body)
	exp, _ := regexp.Compile(`nal":"(.+?)"}`)
	contentArr := exp.FindStringSubmatch(string(buf))
	if len(contentArr) > 1 {
		content = contentArr[1]
	} else {
		_struct.Response(c, errno.ErrExp, contentArr)
	}
	exp, _ = regexp.Compile(`\\`)
	src := exp.ReplaceAllString(content, "")

	var createDate string
	exp, _ = regexp.Compile(`createDate":"(.+?)",`)
	createDateArr := exp.FindStringSubmatch(string(buf))
	if len(createDateArr) > 1 {
		createDate = createDateArr[1]
	} else {
		_struct.Response(c, errno.ErrExp, createDateArr)
	}
	exp, _ = regexp.Compile(`T`)
	createDate1 := exp.ReplaceAllString(createDate, " ")
	exp, _ = regexp.Compile(`\+.+`)
	date := exp.ReplaceAllString(createDate1, "")
	timestamp, _ := time.Parse("2006-01-02 15:04:05", date)
	GMTtime := timestamp.Format("Mon, 02 Jan 2006 15:04:05 GMT")

	imgreq, _ := http.NewRequest("GET", src, nil)
	imgreq.Header.Set("Accept",accept)
	imgreq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	imgreq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7")
	imgreq.Header.Set("Referer", img.Url)
	imgreq.Header.Set("pragma", "no-cache")
	imgreq.Header.Set("Cache-Control", "no-cache")
	imgreq.Header.Set("User-Agent", userAgent)
	imgreq.Header.Set("Upgrade-Insecure-Requests", "1")
	imgreq.Header.Set("If-Modified-Since", GMTtime)
	imgRes, err := client.Do(imgreq)
	fmt.Println(imgRes.Header)
	fmt.Println(img)
	if _, ok := imgreq.Header["Content-Length"]; ok {
		if err != nil || imgRes == nil {
			_struct.Response(c, errno.ErrCurl.Add("imgres"), img.Title)
		}
		defer imgRes.Body.Close()

		imgBytes, err := ioutil.ReadAll(imgRes.Body)
		if err != nil {
			_struct.Response(c, errno.ErrIoutilReadAll, err)
		}

		var suffix string
		exp, _ = regexp.Compile(`p0(.+)`)
		suffixArr := exp.FindStringSubmatch(src)
		if len(createDateArr) > 1 {
			suffix = suffixArr[1]
		} else {
			_struct.Response(c, errno.ErrExp, suffixArr)
		}

		newFile, err := os.Create("static/pixiv/" + img.Title + suffix)
		if err != nil {
			_struct.Response(c, errno.ErrOsCreate, err)
		}
		defer newFile.Close()
		w, err := io.Copy(newFile, bytes.NewReader(imgBytes))
		if w == 0 || err != nil  {
			_struct.Response(c, errno.ErrIoCopy, w)
		}
		fmt.Println(img.Title + suffix + "写入成功")
		// waitGroup.Done()
	} else {
		var suffix string
		exp, _ = regexp.Compile(`p0(.+)`)
		suffixArr := exp.FindStringSubmatch(src)
		if len(createDateArr) > 1 {
			suffix = suffixArr[1]
		} else {
			_struct.Response(c, errno.ErrExp, suffixArr)
		}
		fmt.Println(img.Title + suffix + "写入失败，进入递归")
		fmt.Println(img)
		GetDetail(c, client, img)
	}
}