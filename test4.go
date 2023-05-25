package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"vote/utils"

	"github.com/fatih/color"
	"github.com/google/go-querystring/query"
)

type TestStruct struct {
	name        string
	Token       string
	SendRequest *utils.SendRequest
}

func main() {
	for i := 0; i < 10; i++ {
		test := &TestStruct{
			SendRequest: utils.NewSendRequest(http.Header{}, ""),
		}

	}
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

// 投票
func (selfs *TestStruct) vote(id int) error {
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
		return errors.New(string(responseBody))
	}
	color.Green("o(*￣▽￣*)ブ 投票 %d号 成功 o(*￣▽￣*)ブ \r\n", id)
	return nil
}
