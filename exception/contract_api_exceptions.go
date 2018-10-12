package exception

type ContractApiException struct{ logMessage }

func (ContractApiException) ChainExceptions()       {}
func (ContractApiException) ContractApiExceptions() {}
func (ContractApiException) Code() ExcTypes         { return 3230000 }
func (ContractApiException) What() string           { return "Contract API exception" }

type CryptoApiException struct{ logMessage }

func (CryptoApiException) ChainExceptions()       {}
func (CryptoApiException) ContractApiExceptions() {}
func (CryptoApiException) Code() ExcTypes         { return 3230001 }
func (CryptoApiException) What() string           { return "Crypto API exception" }

type DbApiException struct{ logMessage }

func (DbApiException) ChainExceptions()       {}
func (DbApiException) ContractApiExceptions() {}
func (DbApiException) Code() ExcTypes         { return 3230002 }
func (DbApiException) What() string           { return "Database API exception" }

type ArithmeticException struct{ logMessage }

func (ArithmeticException) ChainExceptions()       {}
func (ArithmeticException) ContractApiExceptions() {}
func (ArithmeticException) Code() ExcTypes         { return 3230003 }
func (ArithmeticException) What() string           { return "Arithmetic exception" }
