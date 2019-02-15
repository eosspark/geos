package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

var BaseUrl string

type client struct {
	HttpClient *http.Client
	BaseURL    string
	Debug      bool
	log        log.Logger
}

func NewClient(baseURL string) *client {
	client := &client{
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
		Debug:   false,
	}
	client.log = log.New("http")
	client.log.SetHandler(log.TerminalHandler)
	return client
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

func (client *client) call(path string, body interface{}) ([]byte, error) {
	jsonBody, err := enc(body)
	if err != nil {
		return nil, err
	}
	targetURL := client.BaseURL + path
	req, err := http.NewRequest("POST", targetURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("NewRequest: %s", err)
	}

	if client.Debug {
		// Useful when debugging API calls
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		client.log.Debug("-------------------------------")
		client.log.Debug(string(requestDump))
		client.log.Debug("")
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", req.URL.String(), err)
	}
	defer resp.Body.Close()

	var cnt bytes.Buffer
	_, err = io.Copy(&cnt, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Copy: %s", err)
	}

	if client.Debug {
		client.log.Debug("Response body: %s", cnt.String())
	}

	statusCode := resp.StatusCode
	if statusCode == 200 || statusCode == 201 || statusCode == 202 {
		return cnt.Bytes(), nil
	} else if statusCode == 404 {
		//Unknown endpoint
		if strings.Contains(path, common.ChainFuncBase) {
			return nil, fmt.Errorf("%s: %s", "Missing Chain API Plugin", targetURL)
		} else if strings.Contains(path, common.WalletFuncBase) {
			return nil, fmt.Errorf("%s: %s", "Missing Wallet API Plugin", targetURL)
		} else if strings.Contains(path, common.HistoryFuncBase) {
			return nil, fmt.Errorf("%s: %s", "Missing History API Plugin", targetURL)
		} else if strings.Contains(path, common.NetFuncBase) {
			return nil, fmt.Errorf("%s: %s", "Missing Net API Plugin", targetURL)
		}
	} else {
		var errorInfo http_plugin.ErrorResults
		err := json.Unmarshal(cnt.Bytes(), &errorInfo)
		if err != nil {
			fmt.Println(err)
		}
		//api.log.Debug("error :%v",errorInfo)
		return nil, fmt.Errorf("Error: %d %s", errorInfo.Error.Code, errorInfo.Error.What)
	}

	if statusCode != 200 {
		client.log.Error("http request fail: Error code %d\n: %s", cnt.String())
		return nil, fmt.Errorf("%s", "http request fail")
	}
	return cnt.Bytes(), nil
}

func DoHttpCall(result interface{}, path string, body interface{}) error {
	client := NewClient(BaseUrl)
	out, err := client.call(path, body)
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
