package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	"vote/config"
	"vote/utils"

	"github.com/fatih/color"
	"github.com/google/go-querystring/query"
	"github.com/henson/proxypool/pkg/models"
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
	directory := "users"
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
	logFile, err := os.OpenFile(fmt.Sprintf("./users/%s.txt", time.Now().Local().Format("2006-01-02")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(logFile, "[info]", log.Ltime)
}

func main() {
	// 初始化配置
	config.InitLog()
	//接收用户的选择
	var id int
	var num int
	var gap int
	var pwd string
	var autoRegister int
	var userPath string
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

	fmt.Printf("是否自动注册，默认不自动注册，自动注册请输入“ 1 ”：")
	fmt.Scanf("%d\n", &autoRegister)
	var userData []map[string]string
	if autoRegister == 0 {
		fmt.Printf("请输入账号本文件路径，默认路径是./user.txt，使用默认路径直接回车：")
		fmt.Scanf("%s\n", &userPath)
		if userPath == "" {
			userPath = "./user.txt"
		}

		file, err := os.Open(userPath)
		if err != nil {
			fmt.Println("无法打开文件:", err)
			return
		}
		defer file.Close()

		// 创建一个Scanner来读取文件内容
		scanner := bufio.NewScanner(file)

		// 逐行读取文件内容
		for scanner.Scan() {
			line := scanner.Text()
			// 账号处理
			// fmt.Println(utils.ParseLogEntry(line))
			userData = append(userData, utils.ParseLogEntry(line))
		}

		// 检查扫描过程是否有错误
		if err := scanner.Err(); err != nil {
			fmt.Println("读取文件错误:", err)
		}

	} else {
		fmt.Printf("密码是否随机默认是，不随机请输入你想要的密码，想随机直接回车：")
		fmt.Scanf("%s\n", &pwd)
	}
	ipChan := make(chan *models.IP, 2000)
	go func() {
		// Start getters to scraper IP and put it in channel
		for {
			go ipProxyRun(ipChan)
			time.Sleep(10 * time.Minute)
		}
	}()
	// fmt.Println(userData)
	// time.Sleep(time.Second * 15)
	headers := &http.Header{}
	// 设置请求头
	headers.Set("Accept", "*/*")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Pragma", "no-cache")
	headers.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	headers.Set("Sec-Ch-Ua-Mobile", "?0")
	headers.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Site", "same-origin")
	headers.Set("X-Requested-With", "XMLHttpRequest")
	headers.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	success := 0
	total := 0
	// 开始进行投票
	for success < num {
		if total > 0 {
			time.Sleep(time.Second * time.Duration(gap))
		}
		headers.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryIxwj5hbrEYmmpCOc")
		vote := &Vote{
			SendRequest: utils.NewSendRequest(headers, "----WebKitFormBoundaryIxwj5hbrEYmmpCOc"),
		}
		// vote.SendRequest.SetProxy("socks5://123.180.0.196:2324", "socks5")
		// 获取代理ip
		select {
		case proxyIp := <-ipChan:
			// 设置代理ip
			proxyurl := fmt.Sprintf("%s://%s", proxyIp.Type1, proxyIp.Data)
			fmt.Println("代理ip：", proxyurl)
			vote.SendRequest.SetProxy(proxyurl, "socks")
		default:
		}
		total++
		err := vote.setToken()
		if err != nil {
			color.Red("获取cookie报错咯：%s\r\n", err)
			continue
		}
		// fmt.Println("Token:", vote.Token)
		// fmt.Println("Cookie:", vote.Cookies)
		if autoRegister == 0 {
			// 使用账号本自动登录
			if total-1 > len(userData) {
				color.Red("账号本的账号都使用完了，无可用账户，请使用另外的账号本！！！\r\n")
				return
			}
			user := userData[total-1]
			if len(user) == 0 {
				color.Red("账号本的账号都使用完了，无可用账户，请使用另外的账号本！！！\r\n")
				continue
			}
			err = vote.login(user["email"], user["pwd"])
			if err != nil {
				color.Red("登录失败咯：%s\r\n", err)
				continue
			}
			// 判断账户是否有投票

		} else {
			// 自动注册
			err = vote.register(pwd)
			if err != nil {
				colorPrint.Add(color.FgRed)
				colorPrint.Println("注册账号报错咯:", err)
				continue
			}
		}
		err = vote.vote(id)
		if err != nil {
			colorPrint.Add(color.FgRed)
			colorPrint.Println("投票报错咯:", err)
			continue
		}
		fmt.Println("======本轮投票结束进行下一次投票======")
		success++
	}

	fmt.Printf("投票结束了，60秒后自动关闭窗口，投给 %d 号明星，总共投票次数：%d ，成功投票：%d\r\n", id, total, success)
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
	body, resp, err := selfs.SendRequest.Get("https://9entertainawards.mcot.net/vote", http.Header{})
	if err != nil {
		return err
	}
	compileRegex := regexp.MustCompile("\"_token\" value=\"(.*?)\">")
	matchArr := compileRegex.FindStringSubmatch(string(body))
	token := matchArr[len(matchArr)-1]
	selfs.Token = token
	cookies := resp.Header["Set-Cookie"]
	selfs.setCookie(cookies)
	// fmt.Println("token:", token)
	// fmt.Println("cookie:", cookie)
	return nil
}

// 登录
func (selfs *Vote) login(email string, pwd string) error {
	// 构建请求参数
	opt := map[string]string{
		"_token":   selfs.Token,
		"email":    email,
		"password": pwd,
	}
	data := utils.MapToUrlValue(opt)
	fmt.Println("请求参数:", data)

	// 进行请求
	responseBody, resp, err := selfs.SendRequest.Post("https://9entertainawards.mcot.net/login", data)
	if err != nil {
		return err
	}
	// 接受请求参数
	response := &ResponseParam{}
	json.Unmarshal(responseBody, &response)
	if response.Code != 200 {
		return errors.New(response.Message)
	}
	fmt.Println("登录成功,账号：", email, "密码：", pwd)
	// 重新设置cookie
	selfs.setCookie(resp.Header["Set-Cookie"])
	selfs.Email = email
	selfs.Password = pwd
	return nil
}

// 判断投票次数

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
	// fmt.Println("注册成功,账号：", opt.EMAIL, "密码：", opt.PASSWORD)
	logger.Println(fmt.Sprintf("注册成功,账号:%s,密码:%s", opt.EMAIL, opt.PASSWORD))
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
	config.Log.Info(fmt.Sprintf("账号:%s,密码:%s,投票 %d 号成功", selfs.Email, selfs.Password, id))
	return nil
}

// 扫描代理IP
func ipProxyRun(ipChan chan<- *models.IP) {
	var wg sync.WaitGroup
	funs := []func() []*models.IP{
		// getter.FQDL,  //新代理 404了
		// getter.PZZQZ, //新代理 不稳定都是超时的
		//getter.Data5u,
		//getter.Feiyi,
		//getter.IP66, //need to remove it
		// getter.IP3306, // 不稳定都是超时的
		// getter.KDL,
		//getter.GBJ,	//因为网站限制，无法正常下载数据
		//getter.Xici,
		//getter.XDL,
		//getter.IP181,  // 已经无法使用
		//getter.YDL,	//失效的采集脚本，用作系统容错实验
		// getter.PLP, //need to remove it
		// getter.PLPSSL,
		// getter.IP89,
	}
	for _, f := range funs {
		wg.Add(1)
		go func(f func() []*models.IP) {
			defer func() {
				if r := recover(); r != nil {
					// 在这里处理panic异常
					// fmt.Println("捕获到panic异常:", r)
				}
			}()
			temp := f()
			// log.Println("[run] get into loop", temp)
			for _, v := range temp {
				log.Println("[run] len of ipChan %v", v)
				// if v.Type1 == "https" {
				ipChan <- v
				// }
			}
			wg.Done()
		}(f)
	}
	wg.Wait()
	log.Println("All getters finished.")
}
