package main

import (
	"log"
	"vote/utils"
)

func main() {

	sendRequest := utils.NewSendRequest(nil, "")
	sendRequest.SetProxy("socks5://194.59.170.116:1080", "socks5")
	// https://api.ip.sb/ip
	body, _, err := sendRequest.Get("https://9entertainawards.mcot.net/vote", nil)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
