package exception

import . "github.com/eosspark/eos-go/log"

type ContractApiException struct{ LogMessage }

func (ContractApiException) ChainExceptions()       {}
func (ContractApiException) ContractApiExceptions() {}
func (ContractApiException) Code() ExcTypes         { return 3230000 }
func (ContractApiException) What() string           { return "Contract API exception" }

type CryptoApiException struct{ LogMessage }

func (CryptoApiException) ChainExceptions()       {}
func (CryptoApiException) ContractApiExceptions() {}
func (CryptoApiException) Code() ExcTypes         { return 3230001 }
func (CryptoApiException) What() string           { return "Crypto API exception" }

type DbApiException struct{ LogMessage }

func (DbApiException) ChainExceptions()       {}
func (DbApiException) ContractApiExceptions() {}
func (DbApiException) Code() ExcTypes         { return 3230002 }
func (DbApiException) What() string           { return "Database API exception" }

type ArithmeticException struct{ LogMessage }

func (ArithmeticException) ChainExceptions()       {}
func (ArithmeticException) ContractApiExceptions() {}
func (ArithmeticException) Code() ExcTypes         { return 3230003 }
func (ArithmeticException) What() string           { return "Arithmetic exception" }
