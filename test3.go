package main

import (
	"fmt"
	"net/http"
	"vote/utils"
)

func main() {
	url := "https://9entertainawards.mcot.net/vote"
	headers := &http.Header{}
	headers.Add("authority", "9entertainawards.mcot.net")
	headers.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Add("accept-language", "zh-CN,zh;q=0.9")
	headers.Add("cache-control", "no-cache")
	headers.Add("cookie", "PHPSESSID=jou8ljr3pgjo2hfq53bqjve2ov; XSRF-TOKEN=eyJpdiI6Ik5GWUhrY0JxN2FyM1MxV2NhWWJReVE9PSIsInZhbHVlIjoialFtbCtwSiszRlRSKzhQcGZQeXRqYkU0MFZYNnljVmsvazdnYTJsc3pVRXFIeXJUWC9rTG9oYTlBNU5KTyt4elJEdzhOMzh2dDJoejlNc0J6M0hoVGNqQVJQaGxZb29tT0lsMy9PMTVTR0pCbm9QUlI4WVA0OG1MdnpqeFdCc1EiLCJtYWMiOiI2ZTA0Y2NhZjVmN2NlMjNkNDRlZWUzMGQ4N2VjODIyNzg5ZDYzYzAwYmY3Y2Y1MmQyODU2YWM5NmE0OWQxYzg4IiwidGFnIjoiIn0%3D; laravel_session=eyJpdiI6ImxoMGQ0V01HWCtOVTFlMGtRU1RqZlE9PSIsInZhbHVlIjoiUkloSjBIbUV3S1EweUc3ZmRiaFNVbUpoZE9IVWpQUzdzdkFkNW0xbzJ4UTh6ck41QlM2dnNQNHRUQjRXbWV5Um05akNtUlBpSDFRelFsMTd3eEYxNkFEc0hWZkNBTitxNWF5VDZ1TDYyR1JvdDYwVVdRdDYrS01vRytXUm9aU2oiLCJtYWMiOiJmZWEzOWY3MzE2NzJmZDYwZGNjZjliZWQwMzdjZTAyZDExY2NhMDNjYmQyYjM2OTgzNDJjOWZlNjNiMDE4N2Y0IiwidGFnIjoiIn0%3D")
	headers.Add("pragma", "no-cache")
	headers.Add("referer", "https://9entertainawards.mcot.net/")
	headers.Add("sec-ch-ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-platform", "\"Windows\"")
	headers.Add("sec-fetch-dest", "document")
	headers.Add("sec-fetch-mode", "navigate")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-user", "?1")
	headers.Add("upgrade-insecure-requests", "1")
	headers.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	headers.Set("User-Agent", utils.GenerateUserAgent())
	req := utils.NewSendRequest(headers, "")

	body, res, err := req.Get(url, http.Header{})
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println(body)
	fmt.Println(res.StatusCode)
}
