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
	"vote/getter"
	"vote/utils"

	"github.com/xuri/excelize/v2"

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
	Token         string
	Cookies       []*http.Cookie
	CookieStr     string
	Phpsessid     *http.Cookie
	LaravelSessid *http.Cookie
	Email         string
	Password      string
	SendRequest   *utils.SendRequest
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

func setUserErr(content ...string) {
	filename := "userError.txt" // 目标文件名
	// 打开文件，使用追加模式打开
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// 写入数据到文件
	_, err = file.WriteString(strings.Join(content, ",") + "\n")
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
	fmt.Printf("请输入停顿时间,单位毫秒,默认10000毫秒：")
	fmt.Scanf("%d\n", &gap)
	if gap == 0 {
		gap = 1000
	}

	fmt.Printf("是否自动注册，默认不自动注册，自动注册请输入“ 1 ”：")
	fmt.Scanf("%d\n", &autoRegister)
	var userData []map[string]string
	if autoRegister == 0 {
		fmt.Printf("请输入账号本文件路径，默认路径是./4444.xlsx：")
		fmt.Scanf("%s\n", &userPath)
		if userPath == "" {
			userPath = "./4444.xlsx"
		}
		fmt.Println("读取文件中……")
		// 判断文件
		if strings.Contains(userPath, "xlsx") {
			// 打开 Excel 文件
			f, err := excelize.OpenFile(userPath)
			if err != nil {
				fmt.Println("无法打开文件:", err)
				return
			}
			// 选择要读取的工作表
			sheetName := "Sheet1"
			rows, err := f.GetRows(sheetName)
			if err != nil {
				fmt.Println("无法打开文件:", err)
				return
			}
			// 读取每一行数据
			for _, row := range rows {
				if len(row) < 2 {
					continue
				}
				user := utils.ParseLogEntry(strings.Join(row, ","))
				// fmt.Println(user)
				// 今日已经用过的账号屏蔽
				if utils.IndexOf(user["email"]) {
					continue
				}
				userData = append(userData, user)
			}
		} else {
			// 读取txt文本
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
				user := utils.ParseLogEntry(line)
				// fmt.Println(user)
				userData = append(userData, user)
			}
			// 检查扫描过程是否有错误
			if err := scanner.Err(); err != nil {
				fmt.Println("读取文件错误:", err)
			}
		}

		fmt.Println("读取完毕总共账号：", len(userData), "个")
		// time.Sleep(time.Hour)

	} else {
		fmt.Printf("密码是否随机默认是，不随机请输入你想要的密码，想随机直接回车：")
		fmt.Scanf("%s\n", &pwd)
	}

	if gap == 0 {
		gap = 10
	}

	ipChan := make(chan *models.IP, 1)
	// go func() {
	// 	// Start getters to scraper IP and put it in channel
	// 	for {
	// 		go ipProxyRun(ipChan)
	// 		time.Sleep(10 * time.Minute)
	// 	}
	// }()
	// fmt.Println(userData)
	// time.Sleep(time.Second * 15)
	headers := &http.Header{}
	// 设置请求头
	headers.Add("authority", "9entertainawards.mcot.net")
	headers.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Add("accept-language", "zh-CN,zh;q=0.9")
	headers.Add("cache-control", "no-cache")
	headers.Add("pragma", "no-cache")
	// headers.Add("referer", "https://9entertainawards.mcot.net/")
	headers.Add("sec-ch-ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"Windows\"")
	headers.Add("sec-fetch-dest", "document")
	headers.Add("sec-fetch-mode", "navigate")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-user", "?1")
	headers.Add("upgrade-insecure-requests", "1")
	headers.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	success := 0
	total := 0
	var wg sync.WaitGroup
	// 开始进行投票
	for success < num || (autoRegister == 0 && success == len(userData)) {
		if total > 0 {
			time.Sleep(time.Millisecond * time.Duration(gap))
		}
		// 算了因为我是分钟收费ip，特殊情况特殊处理一下……
		// go ipProxyRun(ipChan)
		headers.Set("User-Agent", utils.GenerateUserAgent())
		boundary, _ := utils.GenerateRandomBoundary()
		headers.Set("Content-Type", "multipart/form-data; boundary="+boundary)
		vote := &Vote{
			SendRequest: utils.NewSendRequest(headers, boundary),
		}
		// 获取代理ip
		select {
		case proxyIp := <-ipChan:
			// 设置代理ip
			proxyurl := fmt.Sprintf("%s://%s", proxyIp.Type1, proxyIp.Data)
			fmt.Println("代理ip：", proxyurl)
			vote.SendRequest.SetProxy(proxyurl, proxyIp.Type1)
		default:
		}
		total++
		wg.Add(1)
		// 开启协程
		type ft func(int, *Vote, ft)
		f := func(total int, vote *Vote, f ft) {
			err := vote.setToken()
			var user map[string]string
			if autoRegister == 0 && total-1 <= len(userData) {
				user = userData[total-1]
			}
			if err != nil {
				color.Red("投票情况，总投：%d，获取cookie报错咯：%s\r\n", success, err)
				if user != nil {
					time.Sleep(time.Second * 10)
					f(total, &Vote{
						SendRequest: utils.NewSendRequest(headers, boundary),
					}, f)
				}
				return
			}

			if autoRegister == 0 {
				// 使用账号本自动登录
				if total-1 > len(userData) {
					color.Red("账号本的账号都使用完了，无可用账户，请使用另外的账号本！！！\r\n")
					return
				}
				if len(user) == 0 {
					color.Red("账号本的账号都使用完了，无可用账户，请使用另外的账号本！！！\r\n")
					return
				}
				err = vote.login(user["email"], user["pwd"])
				if err != nil {
					color.Red("投票情况，总投：%d，登录失败咯：%s\r\n", success, err)
					time.Sleep(time.Second * 10)
					f(total, &Vote{
						SendRequest: utils.NewSendRequest(headers, boundary),
					}, f)
					return
				}
				// 判断账户是否有投票

			} else {
				// 自动注册
				err = vote.register(pwd)
				if err != nil {
					colorPrint.Add(color.FgRed)
					colorPrint.Println("注册账号报错咯:", err)
					return
				}
			}
			err = vote.vote(id)
			if err != nil {
				colorPrint.Add(color.FgRed)
				color.Red("投票情况，总投：%d，投票报错咯：%s\r\n", success, err)
				time.Sleep(time.Second * 10)
				// f(total, &Vote{
				// 	SendRequest: utils.NewSendRequest(headers, "----WebKitFormBoundaryIxwj5hbrEYmmpCOc"),
				// }, f)
				return
			}
			success++
			fmt.Printf("======本轮投票结束进行下一次投票 当前投了 %d 票======\r\n", success)
		}
		// go func(total int, vote *Vote) {
		// 	defer wg.Done()
		// 	f(total, vote, f)
		// }(total, vote)
		f(total, vote, f)
		wg.Done()
	}
	wg.Wait()
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
		phpsessid, err := req.Cookie("PHPSESSID")
		selfs.Phpsessid = phpsessid
		if err != nil {
			selfs.Phpsessid = &http.Cookie{}
		}
	}
	if selfs.LaravelSessid == nil || selfs.LaravelSessid.Value == "" {
		laravel_session, err := req.Cookie("laravel_session")
		selfs.LaravelSessid = laravel_session
		if err != nil {
			selfs.LaravelSessid = &http.Cookie{}
		}
	}
	xsrftoken, _ := req.Cookie("XSRF-TOKEN")
	cookieStr := fmt.Sprintf("PHPSESSID=%s; XSRF-TOKEN=%s; laravel_session=%s", selfs.Phpsessid.Value, xsrftoken.Value, selfs.LaravelSessid.Value)
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

// 设置token
func (selfs *Vote) setToken2() error {
	// 获取token 和 cookie
	body, resp, err := selfs.SendRequest.Get("https://9entertainawards.mcot.net", http.Header{})
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

	selfs.setToken2()

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
	responseBody, _, err := selfs.SendRequest.RepeatPost("https://9entertainawards.mcot.net/vote/vote", data)
	if err != nil {
		return err
	}
	// 接受请求参数
	// fmt.Println(string(responseBody))
	response := &VoteResponseParam{}
	json.Unmarshal(responseBody, &response)
	if !response.Result {
		color.Red("￣へ￣ 投票失败 ￣へ￣ \r\n")
		config.Log.Error(fmt.Sprintf("账号:%s,密码:%s,投票 %d 号失败", selfs.Email, selfs.Password, id))
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
		// getter.PZZQZ, //新代理
		// getter.IP66,  //need to remove it
		// getter.IP89,
		// getter.Geonode,
		getter.Hsk,
		// getter.IP3306,
		// getter.FQDL,  //新代理 404了
		//getter.Data5u,
		// getter.Feiyi,
		// getter.KDL,
		//getter.GBJ,	//因为网站限制，无法正常下载数据
		// getter.Xici,
		//getter.XDL,
		//getter.IP181,  // 已经无法使用
		//getter.YDL,	//失效的采集脚本，用作系统容错实验
		// getter.PLP, //need to remove it
		// getter.PLPSSL,
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
				// log.Println("[run] len of ipChan %v", v)
				// if v.Type1 == "https" {
				// ip验证是否有效
				sendRequest := utils.NewSendRequest(nil, "")
				if v.Type1 == "" {
					v.Type1 = "socks5"
				}
				sendRequest.SetProxy(v.Type1+"://"+v.Data, v.Type1)
				_, _, err := sendRequest.Get("https://myip.top", nil)
				if err != nil && v.Type2 != "" {
					sendRequest.SetProxy(v.Type2+"://"+v.Data, v.Type2)
					v.Type1 = v.Type2
					_, _, err = sendRequest.Get("https://myip.top", nil)
				}
				if err != nil {
					// fmt.Println("不可用：", proxyAddr)
					return
				}
				// fmt.Printf("%s:%s", v.Type1, v.Data)
				ipChan <- v
				// }
			}
			wg.Done()
		}(f)
	}
	wg.Wait()
	// log.Println("All getters finished.")
}
