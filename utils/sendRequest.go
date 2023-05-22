package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type SendRequest struct {
	client   *http.Client
	retryNum int // 重试次数
	headers  *http.Header
	boundary string
}

func NewSendRequest(headers *http.Header, boundary string) *SendRequest {
	if headers == nil {
		headers = &http.Header{}
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return &SendRequest{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryNum: 5,
		headers:  headers,
		boundary: boundary,
	}
}

func (s *SendRequest) SetHeaders(headers map[string]string) {
	// 设置请求头
	for key, value := range headers {
		s.headers.Set(key, value)
	}
}

func (s *SendRequest) SetProxy(proxyAddr string, t string) {
	// 创建代理 URL
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		fmt.Println("Failed to parse proxy URL:", err)
		return
	}
	var transport *http.Transport
	if t == "socks5" {
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			fmt.Println("Failed to create proxy dialer:", err)
			return
		}

		// 创建自定义的 HTTP 客户端，使用 SOCKS5 代理进行请求
		transport = &http.Transport{
			Dial: dialer.Dial,
		}

		// proxy := func(_ *http.Request) (*url.URL, error) {
		// 	return url.Parse(proxyAddr)
		// }
		// transport = &http.Transport{
		// 	Proxy: proxy,
		// }
	} else {
		// 创建自定义的 Transport
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	s.client.Transport = transport
}

func (s *SendRequest) send(method string, url string, param url.Values, headers http.Header) ([]byte, *http.Response, error) {
	reqBody := strings.NewReader(param.Encode())
	// 设置请求参数
	if s.boundary != "" {
		// 创建请求体
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)
		// 设置分割符号（boundary）
		writer.SetBoundary(s.boundary)
		// 添加表单字段到请求体
		for key, value := range param {
			for _, v := range value {
				_ = writer.WriteField(key, v)
			}
		}
		// 关闭 multipart.Writer，以写入结尾标识符
		_ = writer.Close()
		reqBody = strings.NewReader(buf.String())
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, err
	}

	// 设置请求头
	if headers != nil {
		req.Header = headers
	} else {
		req.Header = *s.headers
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New(fmt.Sprintf("状态码：%d，内容：%s", resp.StatusCode, string(body)))
	}

	return body, resp, nil
}

func (s *SendRequest) Post(url string, param url.Values) ([]byte, *http.Response, error) {
	return s.send("POST", url, param, nil)
}

func (s *SendRequest) RepeatPost(uuid, url string, param url.Values, timeout time.Duration, num int, headers map[string]string) {
	time.Sleep(timeout)

	fmt.Printf("[%s][%d]请求地址：%s\n", uuid, num, url)
	fmt.Printf("[%s][%d]本次发送：%v\n", uuid, num, param)

	result, _, err := s.Post(url, param)
	if err != nil {
		fmt.Printf("[%s][%d]请求返回错误：%v\n", uuid, num, err)
		return
	}
	fmt.Printf("[%s][%d]请求返回：%s\n", uuid, num, result)

	if string(result) != "success" && num < 5 {
		// 进行重发
		s.RepeatPost(uuid, url, param, timeout*2, num+1, headers)
	}
}

func (s *SendRequest) Get(url string, headers http.Header) ([]byte, *http.Response, error) {
	return s.send("GET", url, nil, headers)
}

// func main() {
// 	req := NewSendRequest()

// 	// 测试 POST 请求
// 	postURL := "http://example.com/api"
// 	postParam := url.Values{
// 		"key1": []string{"value1"},
// 		"key2": []string{"value2"},
// 	}
// 	postHeaders := map[string]string{
// 		"Content-Type": "application/x-www-form-urlencoded",
// 	}
// 	postResult, err := req.Post(postURL, postParam, postHeaders)
// 	if err != nil {
// 		fmt.Println("POST 请求错误:", err)
// 	} else {
// 		fmt.Println("POST 请求结果:", postResult)
// 	}

// 	// 测试重复发送 POST 请求
// 	repeatUUID := "12345"
// 	repeatURL := "http://example.com/api"
// 	repeatParam := url.Values{
// 		"key1": []string{"value1"},
// 		"key2": []string{"value2"},
// 	}
// 	repeatTimeout := 0 * time.Second
// 	repeatNum := 1
// 	repeatHeaders := map[string]string{
// 		"Content-Type": "application/x-www-form-urlencoded",
// 	}
// 	req.RepeatPost(repeatUUID, repeatURL, repeatParam, repeatTimeout, repeatNum, repeatHeaders)

// 	// 测试 GET 请求
// 	getURL := "http://example.com/api"
// 	getHeaders := map[string]string{
// 		"Content-Type": "application/x-www-form-urlencoded",
// 	}
// 	getResult, err := req.Get(getURL, getHeaders)
// 	if err != nil {
// 		fmt.Println("GET 请求错误:", err)
// 	} else {
// 		fmt.Println("GET 请求结果:", getResult)
// 	}
// }
