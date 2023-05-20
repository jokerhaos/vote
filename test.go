package main

import (
	"log"
	"vote/utils"
)

func main() {
	// 这个网站有很多代理ip，但是大部分都用不了会超时或者连接拒绝
	// https://geonode.com/free-proxy-list
	sendRequest := utils.NewSendRequest(nil, "")
	sendRequest.SetProxy("socks5://223.215.171.52:2324", "socks5")
	// https://api.ip.sb/ip
	// https://myip.top
	body, _, err := sendRequest.Get("https://myip.top", nil)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
