package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type sendRequest struct {
	client *http.Client
}

func NewSendRequest() *sendRequest {
	return &sendRequest{
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (s *sendRequest) Post(url string, param url.Values, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(param.Encode()))
	if err != nil {
		return "", err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("状态码：%d，内容：%s", resp.StatusCode, string(body)), nil
	}

	return string(body), nil
}

func (s *sendRequest) RepeatPost(uuid, url string, param url.Values, timeout time.Duration, num int, headers map[string]string) {
	time.Sleep(timeout)

	fmt.Printf("[%s][%d]请求地址：%s\n", uuid, num, url)
	fmt.Printf("[%s][%d]本次发送：%v\n", uuid, num, param)

	result, err := s.Post(url, param, headers)
	if err != nil {
		fmt.Printf("[%s][%d]请求返回错误：%v\n", uuid, num, err)
		return
	}
	fmt.Printf("[%s][%d]请求返回：%s\n", uuid, num, result)

	if result != "success" && num < 5 {
		// 进行重发
		s.RepeatPost(uuid, url, param, timeout*2, num+1, headers)
	}
}

func (s *sendRequest) Get(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("状态码：%d，内容：%s", resp.StatusCode, string(body))
	}

	return string(body), nil
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
