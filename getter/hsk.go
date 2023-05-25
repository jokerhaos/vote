package getter

import (
	"net/http"
	"sync"
	"time"
	"vote/utils"

	"github.com/henson/proxypool/pkg/models"
)

var (
	once        sync.Once
	hskSocks5Ip string
)

// GBJ get ip from goubanjia.com
func Hsk() (result []*models.IP) {
	// 执行定时任务，一小时一次
	once.Do(func() {
		// 首次执行任务
		doTask()
		// 创建一个 ticker 定时器，每隔一小时触发一次
		// ticker := time.NewTicker(time.Minute * 58)
		ticker := time.NewTicker(time.Second * 58)
		// 启动一个 goroutine 处理定时任务
		go func() {
			// 循环等待 ticker 定时器的触发事件
			for range ticker.C {
				doTask()
			}
		}()
	})

	ip := models.NewIP()
	ip.Data = hskSocks5Ip
	ip.Type1 = "socks5"
	ip.Source = "huashengdaili"
	result = append(result, ip)
	return
}

func doTask() {
	// pollURL := "https://mobile.huashengdaili.com/servers.php?session=U635687460520130658--b749048e1698763014830e66e5c7230d&time=60&count=1&type=text&pw=no&protocol=s5&ip_type=tunnel"
	pollURL := "https://mobile.huashengdaili.com/servers.php?session=U635687460520130658--b749048e1698763014830e66e5c7230d&time=1&count=1&type=text&pw=no&protocol=s5&ip_type=tunnel"
	req := utils.NewSendRequest(http.Header{}, "")
	body, _, err := req.Get(pollURL, http.Header{})
	if err != nil {
		return
	}
	hskSocks5Ip = string(body)
}
