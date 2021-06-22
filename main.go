package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/evanlixin/go-proxy-examples/handler"
	"github.com/evanlixin/go-proxy-examples/utils"
)

// 参考文档:
// https://blog.csdn.net/zl1zl2zl3/article/details/83149768
// https://developer20.com/writing-proxy-in-go/
// https://zh.wikipedia.org/wiki/X-Forwarded-For
// https://blog.charmes.net/post/reverse-proxy-go/

const (
	defaultServerAddress = "127.0.0.1"
	defaultServerPort    = 8080
)

type ServerOptions struct {
	Address string
	Port    uint
}

func (so *ServerOptions) ValidateServerOptions() error {
	if so == nil {
		return errors.New("serverOptions not init")
	}

	if len(so.Address) == 0 {
		so.Address = defaultServerAddress
	}

	if so.Port == 0 {
		so.Port = defaultServerPort
	}

	return nil
}

var server ServerOptions

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := fmt.Sprintf("%s:%d", server.Address, server.Port)
	fmt.Printf("http Server listen: %s\n", server)

	// register handler
	http.Handle("/", &handler.Proxy{})
	http.HandleFunc("/health", healthHandler)
	go func() {
		err := http.ListenAndServe(server, nil)
		if err != nil {
			fmt.Printf("http ListenAndServe failed: %s\n", err)
			cancel()
		}
	}()

	// register quit
	quitSignalHandler(ctx)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.CommReply(w, r, http.StatusOK, "ok")
}

func init() {
	flag.StringVar(&server.Address, "address", "127.0.0.1", "http server address")
	flag.UintVar(&server.Port, "port", 8080, "http server port")

	flag.Parse()
}

func quitSignalHandler(ctx context.Context) {
	interrupt := make(chan os.Signal, 10)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	// Block until a signal is received.
	select {
	case e := <-interrupt:
		fmt.Println("Got signal:", e.String())
	case <-ctx.Done():
		fmt.Println("Got ctx.Done():", ctx.Err())
	}
}
