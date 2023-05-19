package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 建立与后端服务器的连接
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Println("Failed to connect to destination:", err)
		http.Error(w, "Failed to connect to destination", http.StatusInternalServerError)
		return
	}
	defer destConn.Close()

	// 将客户端的请求发送到后端服务器
	err = r.Write(destConn)
	if err != nil {
		log.Println("Failed to send request to destination:", err)
		http.Error(w, "Failed to send request to destination", http.StatusInternalServerError)
		return
	}

	// 将后端服务器的响应返回给客户端
	_, err = io.Copy(w, destConn)
	if err != nil {
		log.Println("Failed to send response to client:", err)
		http.Error(w, "Failed to send response to client", http.StatusInternalServerError)
		return
	}
}

func main() {
	// 定义处理函数
	// handler := func(w http.ResponseWriter, r *http.Request) {
	// 	// 处理请求
	// 	fmt.Println(r.RemoteAddr)
	// 	fmt.Fprintf(w, "Hello, World!")
	// }

	// // 创建服务器并注册处理函数
	// server := &http.Server{
	// 	Addr:    "0.0.0.0:9090",
	// 	Handler: http.HandlerFunc(handler),
	// }
	// // 启动服务器
	// log.Println("Server started on https://localhost:9090")
	// go func() {
	// 	// if err := server.ListenAndServeTLS("./pem/cert2.pem", "./pem/key2.pem"); err != nil {
	// 	// 	log.Fatal("Server failed to start: ", err)
	// 	// }
	// }()
	// go func() {
	// 	if err := server.ListenAndServe(); err != nil {
	// 		log.Fatal("Server failed to start: ", err)
	// 	}
	// }()

	// 创建代理服务器处理程序
	proxyHandler := http.HandlerFunc(handleRequest)

	// 启动HTTPS代理服务器
	log.Println("proxy 8080")
	err := http.ListenAndServeTLS(":8080", "./pem/cert.pem", "./pem/key.pem", proxyHandler)
	if err != nil {
		log.Fatal("Failed to start proxy server:", err)
	}
}
