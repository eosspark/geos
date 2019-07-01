package fasthttp_test

import (
	"fmt"
	"log"
	"syscall"
	"testing"
	"time"

	"github.com/eosspark/eos-go/libraries/asio"
	"github.com/eosspark/eos-go/plugins/http_plugin/fasthttp"
)

func TestListenAndServe(t *testing.T) {
	listenAddr := "127.0.0.1:8080"

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		time.Sleep(time.Second * 2)
		fmt.Println(time.Now())
		fmt.Fprintf(ctx, "Hello, world! Requested path is %q\n", ctx.Path())
	}

	io := asio.NewIoContext()

	if err := fasthttp.ListenAndAsyncServe(io, listenAddr, requestHandler); err != nil {
		log.Fatalf("error in ListenAndServe: %s\n", err)
	}

	sigint := asio.NewSignalSet(io, syscall.SIGINT)
	sigint.AsyncWait(func(err error) {
		io.Stop()
		sigint.Cancel()
	})

	io.Run()
}

func TestClient(t *testing.T) {
	// Perpare a client, which fetches webpages via HTTP proxy listening
	// on the localhost:8080.
	c := &fasthttp.HostClient{
		Addr: "localhost:8080",
	}

	go func() {
		lastTime, count := time.Now(), 0
		requestHandler := func(ctx *fasthttp.RequestCtx) {
			count++
			if time.Now().Sub(lastTime) >= time.Second {
				fmt.Printf("cont:%d, and %f(us) per request\n", count, float64(1*1e6)/float64(count))
				lastTime = time.Now()
			}
			fmt.Fprintf(ctx, "Hello, world! Requested path is %q\n", ctx.Path())
		}

		io := asio.NewIoContext()

		if err := fasthttp.ListenAndAsyncServe(io, "127.0.0.1:8080", requestHandler); err != nil {
			log.Fatalf("error in ListenAndServe: %s\n", err)
		}

		sigint := asio.NewSignalSet(io, syscall.SIGINT)
		sigint.AsyncWait(func(err error) {
			io.Stop()
			sigint.Cancel()
		})

		io.Run()
	}()

	// Fetch google page via local proxy.
	for {
		statusCode, _, err := c.Get(nil, "http://127.0.0.1:8080/a")
		if err != nil {
			log.Fatalf("Error when loading page through local proxy: %s", err)
		}
		if statusCode != fasthttp.StatusOK {
			log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
	}
	//fmt.Printf("%s\n", body)
	//useResponseBody(body)
}
