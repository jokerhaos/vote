package utils

import (
	"bufio"
	"encoding/base64"
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
	if strings.Contains(logEntry, ":") {
		re := regexp.MustCompile(`账号:(.*?)[，,]*密码:(.*?)($|，|,)`)
		matches := re.FindStringSubmatch(logEntry)
		if len(matches) >= 4 {
			result["email"] = matches[1]
			result["pwd"] = matches[2]
		}
	} else {
		// 分割字符串
		parts := strings.Split(logEntry, ",")
		// 创建map变量并存储账号和密码
		result = map[string]string{
			"email": parts[0],
			"pwd":   parts[1],
		}
	}
	return result
}

func GenerateRandomBoundary() (string, error) {
	// 生成随机的 16 字节
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 将随机字节转换为 base64 编码的字符串
	boundary := base64.RawURLEncoding.EncodeToString(randomBytes)

	// 使用特定前缀，以确保 boundary 的格式符合要求
	boundary = "----WebKitFormBoundary" + boundary

	return boundary, nil
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

func GetHeaders() http.Header {
	// 设置请求头
	headers := http.Header{}
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
	headers.Set("User-Agent", GenerateUserAgent())
	return headers
}
