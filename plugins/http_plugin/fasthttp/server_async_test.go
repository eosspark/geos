package fasthttp_test

import (
	"testing"
	"fmt"
	"log"
	"github.com/eosspark/eos-go/plugins/http_plugin/fasthttp"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"syscall"
)

func TestListenAndServe(t *testing.T) {
	listenAddr := "127.0.0.1:8080"

	requestHandler := func(ctx *fasthttp.RequestCtx) {
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