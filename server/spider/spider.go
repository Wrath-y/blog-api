package spider

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"go-blog/struct"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Img struct {
	ImgId		string
	Title		string
	Url		string
	Praise		string
}

type CountRes struct {
	Success		int
	Failed		int
	List		int
	Page		int
	Exist		int
}

var waitGroup = sync.WaitGroup{}
var maxCh = make(chan int, 4)
var lock = new(sync.Mutex)
var count CountRes
var cookie string

func Get(c *gin.Context, cook string) {
	cookie = cook
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Second * 60,
	}
	GetList(c, client)
	return
}

func GetList(c *gin.Context, client *http.Client)  {
	bookmarkReq, _ := http.NewRequest("GET", bookmark, nil)
	bookmarkReq.Header.Set("cookie", cookie)
	bookmarkResp, err := client.Do(bookmarkReq)
	if err != nil || bookmarkResp == nil {
		_struct.Response(c, errno.ErrCurl.Add("bookmark"), err)
		return
	}
	var buf []byte
	buf, _ = ioutil.ReadAll(bookmarkResp.Body)
	content := string(buf)
	allContent := content
	pageExpInfos, _ := regexp.Compile(`w&amp;p=(\d+)[\s\S]*s="next"`)
	page, _ := strconv.Atoi(pageExpInfos.FindStringSubmatch(content)[1])
	if page == 0 {
		page = 1
	}
	p := 1
	for {
		if p > 1 {
			bookmarkReq, _ = http.NewRequest("GET", bookmark + "?rest=show&p=" + strconv.Itoa(p), nil)
			bookmarkReq.Header.Set("cookie", cookie)
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
	count.Page = p
	count.Success = 0
	count.Failed = 0
	count.Exist = 0
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
		maxCh <- 1
		go GetDetail(c, client, imgSlice[k], false)
		k++
	}
	count.List = k
	go func() {
		waitGroup.Wait()
	}()
	fmt.Println("同步pixiv图片结束")
	_struct.Response(c, nil, count)
	return
}

func GetDetail(c *gin.Context, client *http.Client, img Img, try bool) {
	waitGroup.Add(1)
	defer waitGroup.Done()
	defer func() {
		<- maxCh
	}()
	lock.Lock()
	defer lock.Unlock()
	req, _ := http.NewRequest("GET", img.Url, nil)
	req.Header.Set("cookie", cookie)
	res, err := client.Do(req)
	if err != nil || res == nil {
		fmt.Println(img.Title + img.ImgId + "imgurl")
		//errChan <- err
		//errCodeChan <- errno.ErrCurl.Add(img.Title + img.ImgId + "imgurl")
		count.Failed += 1
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
		fmt.Println(img.Title + img.ImgId + "未匹配到网页内容")
		//errChan <- errno.ErrExp.Add(img.Title + img.ImgId + "未匹配到网页内容")
		//errCodeChan <- errno.ErrExp
		count.Failed += 1
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
		fmt.Println(img.Title + img.ImgId + "未匹配到类型后缀")
		//errChan <- errno.ErrExp.Add(img.Title + img.ImgId + "未匹配到类型后缀")
		//errCodeChan <- errno.ErrExp
		count.Failed += 1
		return
	}

	exp, _ = regexp.Compile(`/`)
	effecTitle := exp.ReplaceAllString(img.Title +  img.ImgId, "-")
	bucket, err := Bucket()
	if err != nil {
		fmt.Println("打开bucket失败")
		//errChan <- err
		//errCodeChan <- errno.UploadError.Add("打开bucket失败")
		count.Failed += 1
		return
	}
	isExist, err := bucket.IsObjectExist(effecTitle + suffix)
	if err != nil {
		fmt.Println(effecTitle + "判断图片是否存在失败")
		//errChan <- err
		//errCodeChan <- errno.UploadError.Add(effecTitle + "判断图片是否存在失败")
		count.Failed += 1
		return
	}
	if isExist == true {
		fmt.Println(effecTitle + suffix + "已存在")
		count.Exist += 1
		return
	}

	imgreq, _ := http.NewRequest("GET", src, nil)
	imgreq.Header.Set("cookie", cookie)
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
			fmt.Println(effecTitle + "未匹配到创建时间")
			//errChan <- errno.ErrExp.Add(effecTitle + "未匹配到创建时间")
			//errCodeChan <- errno.ErrExp
			count.Failed += 1
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
			fmt.Println(effecTitle + "imgres")
			//errChan <- err
			//errCodeChan <- errno.ErrCurl.Add(effecTitle + "imgres")
			count.Failed += 1
			return
		}
		defer imgRes.Body.Close()

		imgBytes, err := ioutil.ReadAll(imgRes.Body)
		if err != nil {
			fmt.Println(effecTitle + "imgBytes")
			//errChan <- err
			//errCodeChan <- errno.ErrIoutilReadAll
			count.Failed += 1
			return
		}

		err = bucket.PutObject(effecTitle + suffix, bytes.NewReader(imgBytes))
		if err != nil {
			fmt.Println(effecTitle + "上传byte数组失败")
			//errChan <- err
			//errCodeChan <- errno.UploadError.Add(effecTitle + "上传byte数组失败")
			count.Failed += 1
			return
		}
		fmt.Println(effecTitle + suffix + "上传成功")
		count.Success += 1
	} else {
		GetDetail(c, client, img, true)
	}

	return
}
