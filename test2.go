package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	logs := []string{
		"[info]13:32:45 账号:f7437cc60d30@gmail.com,密码:f7437cc60d30,投票 3 号成功:",
		"[info]13:32:58 账号:0938e81b039444@gmail.com,密码:0938e81b039444",
		"[info]13:32:58 账号:0938e81b039445@gmail.com，密码:0938e81b039444",
		"[info]13:32:58 账号:0938e81b039446@gmail.com，密码:0938e81b039444，",
		"[info]13:32:58 账号:0938e81b039447@gmail.com,密码:0938e81b039444，",
		"[info]13:32:58 账号:0938e81b039448@gmail.com，密码:0938e81b039444,",
	}

	accounts := make([]map[string]string, 0)

	// regex := regexp.MustCompile(`账号:(.*?)[，,]*密码:(.*?)[，,]*`)
	regex := regexp.MustCompile(`账号:(.*?)[，,]*密码:(.*?)($|，|,)`)

	for _, log := range logs {
		matches := regex.FindStringSubmatch(log)
		if len(matches) == 4 {
			account := make(map[string]string)
			account["email"] = strings.TrimSpace(matches[1])
			account["pwd"] = strings.TrimSpace(matches[2])
			accounts = append(accounts, account)
		}
	}

	for _, account := range accounts {
		fmt.Println(account)
	}
}
