package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"vote/getter"
	"vote/utils"

	"github.com/henson/proxypool/pkg/models"
)

func main() {
	s()
}

type IpData struct {
	Id             string   `json:"_id"`
	Ip             string   `json:"ip"`
	AnonymityLevel string   `json:"anonymityLevel"`
	Asn            string   `json:"asn"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
	Created_at     string   `json:"created_at"`
	Google         string   `json:"google"`
	Isp            string   `json:"isp"`
	LastChecked    string   `json:"lastChecked"`
	Latency        string   `json:"latency"`
	Org            string   `json:"org"`
	Port           string   `json:"port"`
	Protocols      []string `json:"protocols"`
	Region         string   `json:"region"`
	ResponseTime   int      `json:"responseTime"`
	Speed          int      `json:"speed"`
	Updated_at     string   `json:"updated_at"`
}
type IpRequestParam struct {
	Data  []IpData `json:"data"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
}

func s() {

	var wg sync.WaitGroup
	funs := []func() []*models.IP{
		getter.PZZQZ, //新代理
		getter.IP66,  //need to remove it
		getter.IP89,
		getter.Geonode,
		// getter.IP66,  //need to remove it
		// getter.IP89,
		// getter.Geonode,
		// getter.Hsk,
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
				// fmt.Printf("%s:%s\r\n", v.Type1, v.Data)
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
				// }
				setIp(v.Type1 + "://" + v.Data)
			}
			wg.Done()
		}(f)
	}
	wg.Wait()
}

func setIp(content ...string) {
	filename := "ip.txt" // 目标文件名
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
