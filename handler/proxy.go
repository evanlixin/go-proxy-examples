package handler

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

/*
  反向代理
 */

func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		target := targets[rand.Int() % len(targets)]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: director}
}

// 让代理监听在 9090 端口，在后端启动两个返回不同响应的服务器分别监听
// 在 9091 和 9092 端口，通过 curl 访问，可以看到多次请求会返回不同的结果。
// curl http://127.0.0.1:9090
// curl http://127.0.0.1:9090

/*
	proxy := NewMultipleHostsReverseProxy([]*url.URL{
		{
			Scheme: "http",
			Host:   "localhost:9091",
		},
		{
			Scheme: "http",
			Host:   "localhost:9092",
		},
	})
	log.Fatal(http.ListenAndServe(":9090", proxy))
*/