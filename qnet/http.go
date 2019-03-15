package qnet

import (
	"github.com/axgle/mahonia"
	"github.com/camsiabor/qcom/util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type SimpleHttp struct {
	Client *http.Client
}

var simpleHttpInstance = &SimpleHttp{
	Client: &http.Client{},
}

func GetSimpleHttp() *SimpleHttp {

	return simpleHttpInstance
}

func (o *SimpleHttp) SimplePost() {

}

func (o *SimpleHttp) Sleep(millisec int) {
	time.Sleep(time.Millisecond * time.Duration(millisec))
}

func (o *SimpleHttp) Gets(opts []map[string]interface{}) []map[string]interface{} {
	var n = len(opts)
	var waitgroup sync.WaitGroup
	waitgroup.Add(n)
	for i := 0; i < n; i++ {
		go func(one map[string]interface{}) {
			defer waitgroup.Done()
			var url = util.AsStr(one["url"], "")
			var headers = util.AsStringMap(one["headers"], false)
			var encoding = util.AsStr(one["encoding"], "")
			var content, response, err = o.Get(url, headers, encoding)
			one["content"] = content
			one["response"] = response
			one["err"] = err
		}(opts[i])
	}
	waitgroup.Wait()
	return opts
}

func (o *SimpleHttp) Get(url string, headers map[string]string, encoding string) (string, *http.Response, error) {

	var domain string
	var start = strings.Index(url, "://")
	if start < 0 {
		start = 4
		url = "http://" + url
	}
	var fragment = url[start+3:]
	var end = strings.Index(fragment, "/")
	if end < 0 {
		domain = fragment
	} else {
		domain = fragment[:end]
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Host", domain)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	var content string
	resp, err := o.Client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			content = string(bytes[:])
			encoding = strings.ToLower(encoding)
			if encoding != "" && encoding != "utf-8" {
				var encoder = mahonia.NewDecoder(encoding)
				content = encoder.ConvertString(content)
			}
		}
	}

	return content, resp, err
}

func (o *SimpleHttp) Post(url string, headers map[string]string, body string, encoding string) (string, *http.Response, error) {
	var domain string
	var start = strings.Index(url, "://")
	if start < 0 {
		start = 4
		url = "http://" + url
	}
	var fragment = url[start+3:]
	var end = strings.Index(fragment, "/")
	if end < 0 {
		domain = fragment
	} else {
		domain = fragment[:end]
	}

	var bodyio io.Reader
	if len(body) > 0 {
		bodyio = strings.NewReader(body)
	}

	req, err := http.NewRequest("POST", url, bodyio)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Host", domain)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	var content string
	resp, err := o.Client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			content = string(bytes[:])
			encoding = strings.ToLower(encoding)
			if encoding != "" && encoding != "utf-8" {
				var encoder = mahonia.NewDecoder(encoding)
				content = encoder.ConvertString(content)
			}
		}
	}

	return content, resp, err
}
