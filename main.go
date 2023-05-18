package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"vote/utils"

	"github.com/fatih/color"
	"github.com/google/go-querystring/query"
)

type RequestParam struct {
	EMAIL    string `url:"email" json:"email"`
	TOKEN    string `url:"_token" json:"_token"`
	PASSWORD string `url:"password" json:"password"`
}

type ResponseParam struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type VoteRequestParam struct {
	Token        string `url:"_token" json:"_token"`
	Candidate_id int    `url:"candidate_id" json:"candidate_id"`
	Value        int    `url:"value" json:"value"`
	VoteByType   int    `url:"voteByType" json:"voteByType"`
}

type VoteResponseParam struct {
	Reason string `json:"reason"`
	Result bool   `json:"result"`
}

var logger *log.Logger

func init() {
	directory := "logs"
	// 检查目录是否存在
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.Mkdir(directory, 0755) // 设置目录权限
		if err != nil {
			fmt.Println("Failed to create directory:", err)
			return
		}
		fmt.Println("Directory created successfully.")
	} else {
		fmt.Println("Directory already exists.")
	}
	//指定路径的文件，无则创建
	logFile, err := os.OpenFile(fmt.Sprintf("./logs/log_%d.txt", time.Now().Unix()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(logFile, "[info]", log.Ltime)
}
func main() {
	//接收用户的选择
	var id int
	var num int
	var gap int
	var pwd string
	fmt.Println("-----欢迎使用自动投票助手-----")
	fmt.Println("此软件只给M3投票！！！")
	// fmt.Scanf("%d\n", &id)
	id = 3
	fmt.Printf("请输入需要刷票的次数,默认10000000000次：")
	fmt.Scanf("%d\n", &num)
	if num == 0 {
		num = 10000000000
	}
	fmt.Printf("请输入停顿时间,默认10秒：")
	fmt.Scanf("%d\n", &gap)
	if gap == 0 {
		gap = 10
	}

	fmt.Printf("密码是否随机默认是，不随机请输入你想要的密码，想随机直接回车：")
	fmt.Scanf("%s\n", &pwd)

	for i := 0; i < num; i++ {
		vote := &Vote{}
		err := vote.setToken()
		if err != nil {
			fmt.Println("获取cookie报错咯:", err)
			return
		}
		// fmt.Println("Token:", vote.Token)
		// fmt.Println("Cookie:", vote.Cookies)
		err = vote.register(pwd)
		if err != nil {
			fmt.Println("注册账号报错咯:", err)
			return
		}
		err = vote.vote(id)
		if err != nil {
			fmt.Println("投票报错咯:", err)
			return
		}
		fmt.Println("======本轮投票结束进行下一次投票======\r\n")
		time.Sleep(time.Second * time.Duration(gap))
	}

	fmt.Printf("投票结束了，60秒后自动关闭窗口，投给 %d 号明星，总共投票次数：%d \r\n", id, num)
	time.Sleep(time.Second * 60)
}

type Vote struct {
	Token     string
	Cookies   []*http.Cookie
	CookieStr string
	Phpsessid *http.Cookie
	Email     string
	Password  string
}

func (selfs *Vote) setCookie(cookies []string) error {
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", strings.Join(cookies, ";"))
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	if err != nil {
		return err
	}
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	// fmt.Println(strings.Join(cookies, ";"))
	if selfs.Phpsessid == nil || selfs.Phpsessid.Value == "" {
		phpsessid, _ := req.Cookie("PHPSESSID")
		selfs.Phpsessid = phpsessid
	}
	xsrftoken, _ := req.Cookie("XSRF-TOKEN")
	laravel_session, _ := req.Cookie("laravel_session")
	// fmt.Println("===========")
	// fmt.Println(phpsessid.Value)
	// fmt.Println(xsrftoken.Value)
	// fmt.Println(laravel_session.Value)

	cookieStr := fmt.Sprintf("PHPSESSID=%s; XSRF-TOKEN=%s; laravel_session=%s", selfs.Phpsessid.Value, xsrftoken.Value, laravel_session.Value)
	// fmt.Println("cookieStr:", cookieStr)
	selfs.CookieStr = cookieStr
	cookies2 := make([]*http.Cookie, 3)
	cookies2 = append(cookies2, selfs.Phpsessid)
	cookies2 = append(cookies2, xsrftoken)
	cookies2 = append(cookies2, laravel_session)
	selfs.Cookies = cookies2
	return nil
}

// 设置token
func (selfs *Vote) setToken() error {
	// 获取token 和 cookie
	resp, err := http.Get("https://9entertainawards.mcot.net/vote")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	compileRegex := regexp.MustCompile("\"_token\" value=\"(.*?)\">")
	matchArr := compileRegex.FindStringSubmatch(string(body))
	token := matchArr[len(matchArr)-1]
	selfs.Token = token
	cookies := resp.Header["Set-Cookie"]
	selfs.setCookie(cookies)
	// fmt.Println("token:", token)
	// fmt.Println("cookie:", cookie)

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("请求首页状态码异常:%d", resp.StatusCode))
	}
	return nil
}

// 注册
func (selfs *Vote) register(pwd string) error {
	// 构建请求参数
	opt := RequestParam{
		EMAIL:    utils.RandStr2(utils.RandomInt(12, 15)) + "@gmail.com",
		PASSWORD: utils.RandStr2(utils.RandomInt(12, 15)),
		TOKEN:    selfs.Token,
	}
	if pwd != "" {
		opt.PASSWORD = pwd
	}
	data, _ := query.Values(opt)
	fmt.Println("请求参数:", data.Encode())

	// 创建请求体
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	// 设置分割符号（boundary）
	writer.SetBoundary("----WebKitFormBoundaryN2JsHIlOejq9WtWA")

	// 添加表单字段到请求体
	for key, value := range data {
		for _, v := range value {
			_ = writer.WriteField(key, v)
		}
	}
	// 关闭 multipart.Writer，以写入结尾标识符
	_ = writer.Close()

	// reqBody := strings.NewReader("_token=" + selfs.Token + "&email=gfuiasdgkjs16@gmail.com&password=aa123123")
	// boundary := "------WebKitFormBoundaryN2JsHIlOejq9WtWA"
	// reqBody := strings.NewReader("\r\n" + boundary +
	// 	"Content-Disposition: form-data; name=\"email\"\r\n\r\n" +
	// 	"gfuiasdgkjs16@gmail.com\r\n" +
	// 	"------WebKitFormBoundaryN2JsHIlOejq9WtWA\r\n" +
	// 	"Content-Disposition: form-data; name=\"password\"\r\n\r\n" +
	// 	"aa123123\r\n" +
	// 	"------WebKitFormBoundaryN2JsHIlOejq9WtWA\r\n" +
	// 	"Content-Disposition: form-data; name=\"_token\"\r\n\r\n" +
	// 	selfs.Token + "\r\n" +
	// 	"------WebKitFormBoundaryN2JsHIlOejq9WtWA--\r\n")

	// urlValues := url.Values{}
	// urlValues.Add("_token", selfs.Token)
	// urlValues.Add("password", "aa123123")
	// urlValues.Add("email", "gfuiasdgkjs16@gmail.com")
	// fmt.Println("请求参数:", urlValues.Encode())
	// reqBody := strings.NewReader(urlValues.Encode())

	//生成post请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://9entertainawards.mcot.net/register", reqBody)
	if err != nil {
		return err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryN2JsHIlOejq9WtWA")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Android\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://9entertainawards.mcot.net/register")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	req.Header.Set("Cookie", selfs.CookieStr)

	// 设置cookie
	// for _, v := range selfs.Cookies {
	// 	req.AddCookie(v)
	// }

	//Do方法发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("注册状态码异常:%d", resp.StatusCode))
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(responseBody))

	response := &ResponseParam{}
	json.Unmarshal(responseBody, &response)
	if response.Code != 200 {
		return errors.New(response.Message)
	}
	fmt.Println("注册成功,账号：", opt.EMAIL, "密码：", opt.PASSWORD)

	cookies := resp.Header["Set-Cookie"]
	selfs.setCookie(cookies)

	selfs.Email = opt.EMAIL
	selfs.Password = opt.PASSWORD

	return nil
}

// 投票
func (selfs *Vote) vote(id int) error {
	// 构建请求参数
	opt := VoteRequestParam{
		Token:        selfs.Token,
		Candidate_id: id,
		Value:        1,
		VoteByType:   1,
	}
	data, _ := query.Values(opt)
	fmt.Println("请求参数:", data.Encode())

	// 创建请求体
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	// 设置分割符号（boundary）
	writer.SetBoundary("----WebKitFormBoundaryN2JsHIlOejq9WtWA")

	// 添加表单字段到请求体
	for key, value := range data {
		for _, v := range value {
			_ = writer.WriteField(key, v)
		}
	}
	// 关闭 multipart.Writer，以写入结尾标识符
	_ = writer.Close()

	//生成post请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://9entertainawards.mcot.net/vote/vote", reqBody)
	if err != nil {
		return err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryN2JsHIlOejq9WtWA")
	req.Header.Set("Cookie", selfs.CookieStr)

	// 设置cookie
	// for _, v := range selfs.Cookies {
	// 	req.AddCookie(v)
	// }

	//Do方法发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("投票状态码异常:%d", resp.StatusCode))
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(responseBody))

	response := &VoteResponseParam{}
	json.Unmarshal(responseBody, &response)
	colorPrint := color.New()

	if !response.Result {
		colorPrint.Add(color.FgRed) // 红色文字
		colorPrint.Println("￣へ￣ 投票失败 ￣へ￣")
		return errors.New(string(responseBody))
	}
	colorPrint.Add(color.FgGreen)
	colorPrint.Printf("o(*￣▽￣*)ブ 投票 %d号 成功 o(*￣▽￣*)ブ \r\n", id)
	logger.Println(fmt.Sprintf("账号:%s,密码:%s,投票 %d 号成功:", selfs.Email, selfs.Password, id))
	return nil
}
