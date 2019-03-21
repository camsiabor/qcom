package qnet

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/camsiabor/qcom/util"
	gorilla "github.com/gorilla/http"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type SimpleHttp struct {
	stds         []*http.Client
	gorillas     []*gorilla.Client
	stdCount     int
	gorillaCount int
	mutex        sync.RWMutex
}

var simpleHttpInstance = &SimpleHttp{}

func GetSimpleHttp() *SimpleHttp {
	return simpleHttpInstance
}

func (o *SimpleHttp) InitClients(t string, count int, timeout int) {

	if count <= 1 {
		count = 2
	}

	rand.Seed(time.Now().UnixNano())

	o.mutex.Lock()
	defer o.mutex.Unlock()

	if t == "std" {
		o.stdCount = count
		o.stds = make([]*http.Client, count)
		for i := 0; i < o.stdCount; i++ {
			var client = &http.Client{}
			o.stds[i] = client
			client.Timeout = time.Duration(timeout) * time.Second
		}
	} else {
		o.gorillaCount = count
		o.gorillas = make([]*gorilla.Client, count)
		for i := 0; i < o.gorillaCount; i++ {
			var client = &gorilla.Client{}
			o.gorillas[i] = client
		}
	}

}

func (o *SimpleHttp) GetClient(t string) interface{} {
	if o.stds == nil {
		o.InitClients("std", 22, 15)
	}
	if o.gorillas == nil {
		o.InitClients("gorilla", 22, 15)
	}
	var r interface{}
	o.mutex.RLock()
	if t == "std" {
		var index = rand.Int() % o.stdCount
		r = o.stds[index]
	} else {
		var index = rand.Int() % o.gorillaCount
		r = o.gorillas[index]
	}
	o.mutex.RUnlock()
	return r
}

func (o *SimpleHttp) SimplePost() {

}

func (o *SimpleHttp) Sleep(millisec int) {
	time.Sleep(time.Millisecond * time.Duration(millisec))
}

func (o *SimpleHttp) Gets(t string, opts []map[string]interface{}) []map[string]interface{} {
	var n = len(opts)
	var waitgroup sync.WaitGroup
	waitgroup.Add(n)
	for i := 0; i < n; i++ {
		go func(one map[string]interface{}) {
			defer waitgroup.Done()
			var url = util.AsStr(one["url"], "")
			var headers = util.AsStringMap(one["headers"], false)
			var encoding = util.AsStr(one["encoding"], "")
			var content, response, err = o.Get(t, url, headers, encoding)
			one["content"] = content
			one["response"] = response
			one["err"] = err
		}(opts[i])
	}
	waitgroup.Wait()
	return opts
}

func (o *SimpleHttp) Get(t string, url string, headers map[string]string, encoding string) (string, interface{}, error) {
	if t == "std" {
		return o.StdGet(url, headers, encoding)
	} else {
		return o.GorillaGet(url, headers, encoding)
	}
}

func (o *SimpleHttp) GorillaGet(url string, headers map[string]string, encoding string) (string, interface{}, error) {

	var gheaders map[string][]string
	if headers != nil {
		gheaders = make(map[string][]string)
		for k, v := range headers {
			var one = make([]string, 1)
			one[0] = v
			gheaders[k] = one
		}
	}
	var client = gorilla.DefaultClient // o.GetClient("gorilla").(*gorilla.Client)
	var status, respheaders, reader, err = client.Get(url, gheaders)
	if err != nil {
		return "", respheaders, err
	}
	if status.Code != 200 {
		return "", respheaders, fmt.Errorf("response status %v", status)
	}
	defer reader.Close()
	bytes, err := ioutil.ReadAll(reader)
	var content string
	if err == nil {
		content = string(bytes[:])
		encoding = strings.ToLower(encoding)
		if encoding != "" && encoding != "utf-8" {
			var encoder = mahonia.NewDecoder(encoding)
			content = encoder.ConvertString(content)
		}
	}
	return content, respheaders, err
}

func (o *SimpleHttp) StdGet(url string, headers map[string]string, encoding string) (string, interface{}, error) {

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
	var client = o.GetClient("std").(*http.Client)
	resp, err := client.Do(req)
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
	var client = o.GetClient("std").(*http.Client)
	resp, err := client.Do(req)
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
