package utils

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func ParseCookieString(rawCookies string) (*http.Request, error) {
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", rawCookies)
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	return req, err
}

// var utilsBytesSeed []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

func init() {
	// 保证每次生成的随机数不一样
	rand.Seed(time.Now().UnixNano())
}

// 方法二
func RandStr2(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}

func RandomInt(min, max int) int {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机数
	return rand.Intn(max-min+1) + min
}

// map转url.Values
func MapToUrlValue(params map[string]string) url.Values {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	return values
}

func ParseLogEntry(logEntry string) map[string]string {
	result := make(map[string]string)

	// 提取账号和密码信息
	// re := regexp.MustCompile(`账号:(.*?),密码:(.*?),`)
	re := regexp.MustCompile(`账号:(.*?)[，,]*密码:(.*?)($|，|,)`)
	matches := re.FindStringSubmatch(logEntry)
	if len(matches) >= 4 {
		result["email"] = matches[1]
		result["pwd"] = matches[2]
	}

	return result
}

func ParseLogEntry2(logEntry []string) map[string]string {
	result := make(map[string]string)

	// 提取账号和密码信息
	// re := regexp.MustCompile(`账号:(.*?),密码:(.*?),`)
	re := regexp.MustCompile(`账号:(.*?)[，,]*密码:(.*?)($|，|,)`)
	matches := re.FindStringSubmatch(strings.Join(logEntry, ","))
	result["email"] = matches[1]
	result["pwd"] = matches[2]

	return result
}

func IndexOf(targetString string) bool {
	filePath := fmt.Sprintf("./logs/info-%s.log", time.Now().Local().Format("2006-01-02"))
	// 读取文件内容
	content, _ := ioutil.ReadFile(filePath)
	filePath2 := fmt.Sprintf("./logs/error-%s.log", time.Now().Local().Format("2006-01-02"))
	// 读取文件内容
	content2, _ := ioutil.ReadFile(filePath2)
	// 将文件内容转换为字符串
	fileContent := string(content) + string(content2)
	// 判断目标字符串是否在文件内容中
	return strings.Contains(fileContent, targetString)
}
