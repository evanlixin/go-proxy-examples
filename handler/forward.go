package handler

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

/*
正向代理
 */

type Proxy struct{}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	transport :=  http.DefaultTransport

	// step 1 代理接受到客户端的请求，复制原来的请求对象，并根据数据配置新请求的各种参数(添加上X-Forward-For头部等)
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		prior, ok := outReq.Header["X-Forwarded-For"]
		if ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step 2 把新请求复制到服务器端，并接收到服务器端返回的响应
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	// step 3 代理服务器对响应做一些处理，然后返回给客户端
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
}

// 代码运行之后，会在本地的 8080 端口启动代理服务。修改浏览器的代理为 127.0.0.1:8080
// 再访问网站，可以验证代理正常工作，也能看到它在终端打印出所有的请求信息。
