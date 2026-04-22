package main

import (
	"fmt"
	"strings"

	gamemitm "github.com/wangggeng/game-mitm"
	"github.com/wangggeng/game-mitm/gosysproxy"

	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := gosysproxy.SetGlobalProxy(
		"127.0.0.1:12311",
		"localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;192.168.*",
	)
	if err != nil {
		panic(err)
	}
}
func main() {
	proxy := gamemitm.NewProxy()
	proxy.SetVerbose(false)
	proxy.SetLogger(gamemitm.NewDefaultLogger(int(gamemitm.FATAL) + 1))

	// 拦截所有HTTPS请求
	proxy.OnRequest("*").Do(func(body []byte, ctx *gamemitm.ProxyCtx) []byte {
		return body
	})

	// 拦截指定接口的响应
	proxy.OnResponse("*").Do(func(body []byte, ctx *gamemitm.ProxyCtx) []byte {
		fullURL := ctx.Req.Host + ctx.Req.URL.RequestURI()
		// 只打印匹配的接口
		if strings.Contains(fullURL, "jy.dxywjy.cn") {
			fmt.Println("========== 命中目标接口 ==========")
			fmt.Printf("URL: https://%s\n", fullURL)
			fmt.Printf("状态码: %d\n", ctx.Resp.StatusCode)
			fmt.Printf("Content-Type: %s\n", ctx.Resp.Header.Get("Content-Type"))
			fmt.Printf("响应体(%d bytes):\n%s\n", len(body), string(body))
			fmt.Printf("请求头：%s\n", ctx.Req.Header)
			fmt.Printf("响应头：%s\n", ctx.Resp.Header)
			fmt.Println("==================================")
		}
		return body
	})

	// 监听操作系统信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动代理
	go proxy.Start()

	// 等待程序终止信号
	<-signalChan
	gosysproxy.Off()
	// 在程序结束时执行清理操作

}

func truncate(data []byte, maxLen int) string {
	s := string(data)
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}
