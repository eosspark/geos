package wallet_plugin

import "github.com/eosspark/eos-go/crypto/ecc"

type BaseWalletApi interface {
	GetPrivateKey(pubkey ecc.PublicKey) ecc.PrivateKey

	IsLocked() bool
	Lock()
	Unlock(password string)

	CheckPassword(password string)
	SetPassword(password string)

	ListKeys() map[ecc.PublicKey]ecc.PrivateKey
	ListPublicKeys() []ecc.PublicKey

	ImportKey(wifKey string) bool
	RemoveKey(key string) bool
	CreateKey(keyType string) string

	TrySignDigest(digest []byte, publicKey ecc.PublicKey) *ecc.Signature
}