package wallet_plugin

import (
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
)

type KeyMapType = map[ecc.PublicKey]uint16

type YhConnector struct {
}

type YhSession struct {
}

type YhCapabilities struct {
}

type YubihsmApi struct {
}

type YubihsmWalletImpl struct {
	Connector      *YhConnector
	Session        *YhSession
	Endpoint       string
	AuthKey        uint16
	Keys           KeyMapType
	AuthKeyCaps    YhCapabilities
	AuthKeyDomains uint16
	Api            YubihsmApi
}

func (y YubihsmWalletImpl) IsLocked() bool {
	return y.Connector == nil
}

func (y YubihsmWalletImpl) Unlock(password string) {

}

func (y YubihsmWalletImpl) Lock() {
	if y.Session != nil {

	}
}

type YubihsmWallet struct {
	my *YubihsmWalletImpl
}

func (y YubihsmWallet) GetPrivateKey(pubkey ecc.PublicKey) ecc.PrivateKey {
	try.EosThrow(&exception.WalletException{}, "Obtaining private key for a key stored in YubiHSM is impossible")
	return ecc.PrivateKey{}
}

func (y YubihsmWallet) IsLocked() bool {
	return y.my.IsLocked()
}

func (y YubihsmWallet) Lock() {
	try.FcAssert(!y.IsLocked())
	y.my.Lock()
}
