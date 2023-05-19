package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

type Vote struct {
	Token       string
	Cookies     []*http.Cookie
	CookieStr   string
	Phpsessid   *http.Cookie
	Email       string
	Password    string
	SendRequest *utils.SendRequest
}

var logger *log.Logger
var colorPrint = color.New()

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

	headers := &http.Header{}
	// 设置请求头
	headers.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryN2JsHIlOejq9WtWA")
	headers.Set("Accept", "*/*")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Pragma", "no-cache")
	headers.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	headers.Set("Sec-Ch-Ua-Mobile", "?1")
	headers.Set("Sec-Ch-Ua-Platform", "\"Android\"")
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Site", "same-origin")
	headers.Set("X-Requested-With", "XMLHttpRequest")
	headers.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	// 开始进行投票
	for i := 0; i < num; i++ {
		vote := &Vote{
			SendRequest: utils.NewSendRequest(headers, "----WebKitFormBoundaryN2JsHIlOejq9WtWA"),
		}
		err := vote.setToken()
		if err != nil {
			colorPrint.Add(color.FgRed)
			colorPrint.Println("获取cookie报错咯:", err)
			continue
		}
		// fmt.Println("Token:", vote.Token)
		// fmt.Println("Cookie:", vote.Cookies)
		err = vote.register(pwd)
		if err != nil {
			colorPrint.Add(color.FgRed)
			colorPrint.Println("注册账号报错咯:", err)
			continue
		}
		err = vote.vote(id)
		if err != nil {
			colorPrint.Add(color.FgRed)
			colorPrint.Println("投票报错咯:", err)
			continue
		}
		fmt.Println("======本轮投票结束进行下一次投票======")
		time.Sleep(time.Second * time.Duration(gap))
	}

	fmt.Printf("投票结束了，60秒后自动关闭窗口，投给 %d 号明星，总共投票次数：%d \r\n", id, num)
	time.Sleep(time.Second * 60)
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
	cookieStr := fmt.Sprintf("PHPSESSID=%s; XSRF-TOKEN=%s; laravel_session=%s", selfs.Phpsessid.Value, xsrftoken.Value, laravel_session.Value)
	selfs.SendRequest.SetHeaders(map[string]string{
		"Cookie": cookieStr,
	})
	// fmt.Println("cookieStr:", cookieStr)
	// selfs.CookieStr = cookieStr
	// cookies2 := make([]*http.Cookie, 3)
	// cookies2 = append(cookies2, selfs.Phpsessid)
	// cookies2 = append(cookies2, xsrftoken)
	// cookies2 = append(cookies2, laravel_session)
	// selfs.Cookies = cookies2
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
	// 进行请求
	responseBody, resp, err := selfs.SendRequest.Post("https://9entertainawards.mcot.net/register", data)
	if err != nil {
		return err
	}
	// 接受请求参数
	response := &ResponseParam{}
	json.Unmarshal(responseBody, &response)
	if response.Code != 200 {
		return errors.New(response.Message)
	}
	fmt.Println("注册成功,账号：", opt.EMAIL, "密码：", opt.PASSWORD)
	// 重新设置cookie
	selfs.setCookie(resp.Header["Set-Cookie"])
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

	// 进行请求
	responseBody, _, err := selfs.SendRequest.Post("https://9entertainawards.mcot.net/vote/vote", data)
	if err != nil {
		return err
	}
	// 接受请求参数
	// fmt.Println(string(responseBody))
	response := &VoteResponseParam{}
	json.Unmarshal(responseBody, &response)
	if !response.Result {
		color.Red("￣へ￣ 投票失败 ￣へ￣ \r\n")
		return errors.New(string(responseBody))
	}
	color.Green("o(*￣▽￣*)ブ 投票 %d号 成功 o(*￣▽￣*)ブ \r\n", id)
	logger.Println(fmt.Sprintf("账号:%s,密码:%s,投票 %d 号成功:", selfs.Email, selfs.Password, id))
	return nil
}
