package utils

import (
	"fmt"
	"math/rand"
	"time"
)

type UserAgent struct {
	OS      string
	Browser string
}

func getRandomOS() string {
	// 随机选择操作系统
	osList := []string{
		"Windows NT 10.0",
		"Macintosh; Intel Mac OS X 10_15_7",
		"X11; Linux x86_64",
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(osList))
	return osList[index]
}

func getRandomBrowser() string {
	// 随机选择浏览器
	browserList := []string{
		"Chrome/94.0.4606.71",
		"Safari/537.36",
		"Firefox/93.0",
		"Edge/94.0.992.38",
		"Opera/78.0.4093.184",
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(browserList))
	return browserList[index]
}

func GenerateUserAgent() string {
	os := getRandomOS()
	browser := getRandomBrowser()
	userAgent := fmt.Sprintf("Mozilla/5.0 (%s; %s) AppleWebKit/537.36 (KHTML, like Gecko) %s", os, browser, browser)
	return userAgent
}
