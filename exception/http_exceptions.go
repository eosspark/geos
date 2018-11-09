package exception

import . "github.com/eosspark/eos-go/log"

type HttpException struct{ LogMessage }

func (HttpException) ChainExceptions() {}
func (HttpException) HttpExceptions()  {}
func (HttpException) Code() ExcTypes   { return 3200000 }
func (HttpException) What() string {
	return "http exception"
}

type InvalidHttpClientRootCert struct{ LogMessage }

func (InvalidHttpClientRootCert) ChainExceptions() {}
func (InvalidHttpClientRootCert) HttpExceptions()  {}
func (InvalidHttpClientRootCert) Code() ExcTypes   { return 3200001 }
func (InvalidHttpClientRootCert) What() string {
	return "invalid http client root certificate"
}

type InvalidHttpResponse struct{ LogMessage }

func (InvalidHttpResponse) ChainExceptions() {}
func (InvalidHttpResponse) HttpExceptions()  {}
func (InvalidHttpResponse) Code() ExcTypes   { return 3200002 }
func (InvalidHttpResponse) What() string {
	return "invalid http response"
}

type ResolvedToMultiplePorts struct{ LogMessage }

func (ResolvedToMultiplePorts) ChainExceptions() {}
func (ResolvedToMultiplePorts) HttpExceptions()  {}
func (ResolvedToMultiplePorts) Code() ExcTypes   { return 3200003 }
func (ResolvedToMultiplePorts) What() string {
	return "service resolved to multiple ports"
}

type FailToResolveHost struct{ LogMessage }

func (FailToResolveHost) ChainExceptions() {}
func (FailToResolveHost) HttpExceptions()  {}
func (FailToResolveHost) Code() ExcTypes   { return 3200004 }
func (FailToResolveHost) What() string {
	return "fail to resolve host"
}

type HttpRequestFail struct{ LogMessage }

func (HttpRequestFail) ChainExceptions() {}
func (HttpRequestFail) HttpExceptions()  {}
func (HttpRequestFail) Code() ExcTypes   { return 3200005 }
func (HttpRequestFail) What() string {
	return "http request fail"
}

type InvalidHttpRequest struct{ LogMessage }

func (InvalidHttpRequest) ChainExceptions() {}
func (InvalidHttpRequest) HttpExceptions()  {}
func (InvalidHttpRequest) Code() ExcTypes   { return 3200006 }
func (InvalidHttpRequest) What() string {
	return "invalid http request"
}
