package spider

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func Login() {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	value := url.Values{}
	value.Add("email", email)
	value.Add("password", password)
	body := ioutil.NopCloser(strings.NewReader(value.Encode()))
	client.Post(login_url, "application/x-www-form-urlencoded; charset=UTF-8", body)
}