package console

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/log"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

var ErrNotFound = errors.New("resource not found")

type API struct {
	HttpClient *http.Client
	BaseURL    string
	Debug      bool
	//Compress                common.CompressionType
	DefaultMaxCPUUsageMS    uint8
	DefaultMaxNetUsageWords uint32 // in 8-bytes words
	log                     log.Logger
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
		BaseURL: baseURL,
		//Compress: common.CompressionZlib,
		Debug: false,
	}
	api.log = log.New("http")
	api.log.SetHandler(log.TerminalHandler)
	return api
}

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

func (api *API) call(path string, body interface{}) ([]byte, error) {
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
		api.log.Debug("-------------------------------")
		api.log.Debug(string(requestDump))
		api.log.Debug("")
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

	//api.log.Debug("Response body: %s", cnt.String())

	if api.Debug {
		api.log.Debug("Response body: %s", cnt.String())
	}

	return cnt.Bytes(), nil
}

func DoHttpCall(result interface{}, path string, body interface{}) error {
	http := NewHttp(common.HttpEndPoint)
	out, err := http.call(path, body)
	if err != nil {
		return err
	}

	if result != nil {
		if err := json.Unmarshal(out, &result); err != nil {
			fmt.Printf("dohttpCall, Unmarshal: %s\n", err.Error())
			return err
		}
	}

	return nil
}
