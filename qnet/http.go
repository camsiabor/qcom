package qnet

import (
	"io/ioutil"
	"net/http"
	"strings"
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

func (o *SimpleHttp) Get(url string, headers map[string]string, encoding string) (string, error) {

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
		return "", err
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
			content = string(bytes)
		}
	}

	return content, err
}
