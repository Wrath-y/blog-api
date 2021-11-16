package spider

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/server/errno"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Img struct {
	ImgId  string
	Title  string
	Url    string
	Praise string
}

type CountRes struct {
	Success int
	Failed  int
	List    int
	Page    int
	Exist   int
}

type Conf struct {
	waitGroup *sync.WaitGroup
	maxCh     chan int
	lock      *sync.Mutex
	count     CountRes
	cookie    string
}

func Get(c *gin.Context, cook string) {
	conf := new(Conf)
	conf.waitGroup = new(sync.WaitGroup)
	conf.maxCh = make(chan int, 100)
	conf.lock = new(sync.Mutex)
	conf.count = CountRes{}
	conf.cookie = cook

	proxy, _ := url.Parse("socks5://127.0.0.1:1080")
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 60,
	}
	GetList(c, client, conf)
	return
}

func GetList(c *gin.Context, client *http.Client, conf *Conf) {
	bookmarkReq, _ := http.NewRequest("GET", bookmark, nil)
	bookmarkReq.Header.Set("cookie", conf.cookie)
	bookmarkResp, err := client.Do(bookmarkReq)
	if err != nil || bookmarkResp == nil {
		controller.Response(c, errno.CurlErr.Add("bookmark"), err)
		return
	}
	var buf []byte
	buf, _ = ioutil.ReadAll(bookmarkResp.Body)
	content := string(buf)
	fmt.Println(content)
	allContent := content
	pageExpInfos, _ := regexp.Compile(`w&amp;p=(\d+)[\s\S]*s="next"`)
	if len(pageExpInfos.FindStringSubmatch(content)) == 0 {
		controller.Response(c, errno.IndexOutOfRangeErr.Add("pageExpInfos.FindStringSubmatch(content)"), err)
		return
	}
	page, err := strconv.Atoi(pageExpInfos.FindStringSubmatch(content)[1])
	if err != nil {
		controller.Response(c, errno.RegexpErr.Add("pageExpInfos.FindStringSubmatch(content)"), err)
		return
	}
	if page == 0 {
		page = 1
	}
	p := 1
	for {
		if p > 1 {
			bookmarkReq, _ = http.NewRequest("GET", bookmark+"?rest=show&p="+strconv.Itoa(p), nil)
			bookmarkReq.Header.Set("cookie", conf.cookie)
			bookmarkResp, err = client.Do(bookmarkReq)
			if bookmarkResp == nil {
				controller.Response(c, errno.CurlErr.Add("bookmarkwithpage"), err)
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
	conf.count.Page = p
	conf.count.Success = 0
	conf.count.Failed = 0
	conf.count.Exist = 0
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
		conf.maxCh <- 1
		go GetDetail(c, client, conf, imgSlice[k], false)
		k++
	}
	conf.count.List = k
	go func() {
		conf.waitGroup.Wait()
	}()
	fmt.Println("同步pixiv图片结束")
	controller.Response(c, nil, conf.count)

	return
}

func GetDetail(c *gin.Context, client *http.Client, conf *Conf, img Img, try bool) {
	conf.waitGroup.Add(1)
	defer conf.waitGroup.Done()
	defer func() {
		<-conf.maxCh
	}()
	conf.lock.Lock()
	defer conf.lock.Unlock()
	req, _ := http.NewRequest("GET", img.Url, nil)
	req.Header.Set("cookie", conf.cookie)
	res, err := client.Do(req)
	if err != nil || res == nil {
		fmt.Println(img.Title + img.ImgId + "imgurl")
		conf.count.Failed += 1
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
		conf.count.Failed += 1
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
		conf.count.Failed += 1
		return
	}

	exp, _ = regexp.Compile(`/`)
	effecTitle := exp.ReplaceAllString(img.Title+img.ImgId, "-")
	bucket, err := Bucket()
	if err != nil {
		fmt.Println("打开bucket失败")
		conf.count.Failed += 1
		return
	}
	isExist, err := bucket.IsObjectExist(effecTitle + suffix)
	if err != nil {
		fmt.Println(effecTitle + "判断图片是否存在失败")
		conf.count.Failed += 1
		return
	}
	if isExist == true {
		fmt.Println(effecTitle + suffix + "已存在")
		conf.count.Exist += 1
		return
	}

	imgreq, _ := http.NewRequest("GET", src, nil)
	imgreq.Header.Set("cookie", conf.cookie)
	imgreq.Header.Set("Accept", accept)
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
			conf.count.Failed += 1
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
			conf.count.Failed += 1
			return
		}
		defer imgRes.Body.Close()

		imgBytes, err := ioutil.ReadAll(imgRes.Body)
		if err != nil {
			fmt.Println(effecTitle + "读取imgBytes失败")
			conf.count.Failed += 1
			return
		}

		err = bucket.PutObject(effecTitle+suffix, bytes.NewReader(imgBytes))
		if err != nil {
			fmt.Println(effecTitle + "上传byte数组失败")
			conf.count.Failed += 1
			return
		}
		fmt.Println(effecTitle + suffix + "上传成功")
		conf.count.Success += 1
	} else {
		GetDetail(c, client, conf, img, true)
	}

	return
}
