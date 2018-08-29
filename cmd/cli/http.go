package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

var ErrNotFound = errors.New("resource not found")

type API struct {
	HttpClient              *http.Client
	BaseURL                 string
	Debug                   bool
	Compress                common.CompressionType
	DefaultMaxCPUUsageMS    uint8
	DefaultMaxNetUsageWords uint32 // in 8-bytes words
}

func NewHttp(baseURL string) *API {
	api := &API{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true, // default behavior, because of `nodeos`'s lack of support for Keep alives.
			},
		},
		BaseURL:  baseURL,
		Compress: common.CompressionZlib,
		// Debug:    true,
	}

	return api
}

// See more here: libraries/chain/contracts/abi_serializer.cpp:58...

func (api *API) call(path string, body interface{}, out interface{}) ([]byte, error) {
	jsonBody, err := enc(body)
	if err != nil {
		return nil, err
	}
	targetURL := api.BaseURL + path
	// targetURL := fmt.Sprintf("%s/v1/%s/%s", api.BaseURL, baseAPI, endpoint)
	req, err := http.NewRequest("POST", targetURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("NewRequest: %s", err)
	}

	if api.Debug {
		// Useful when debugging API calls
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("-------------------------------")
		fmt.Println(string(requestDump))
		fmt.Println("")
	}

	resp, err := api.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", req.URL.String(), err)
	}
	defer resp.Body.Close()

	var cnt bytes.Buffer
	_, err = io.Copy(&cnt, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Copy: %s", err)
	}

	if resp.StatusCode == 404 {
		return nil, ErrNotFound
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("%s: status code=%d, body=%s", req.URL.String(), resp.StatusCode, cnt.String())
	}

	if api.Debug {
		fmt.Println("RESPONSE:")
		fmt.Println(cnt.String())
		fmt.Println("")
	}
	// fmt.Println("返回数据： ", cnt)

	if err := json.Unmarshal(cnt.Bytes(), &out); err != nil {
		return nil, fmt.Errorf("Unmarshal: %s", err)
	}

	return cnt.Bytes(), nil
}

type M map[string]interface{}

func enc(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}

	cnt, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(cnt), nil
}

func DoHttpCall(path string, body interface{}, out interface{}) (data []byte, err error) {
	http := NewHttp("http://127.0.0.1:8888")
	data, err = http.call(path, body, &out)
	return
}
