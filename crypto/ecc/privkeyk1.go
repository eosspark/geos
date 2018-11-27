package ecc

import (
	"fmt"

	"github.com/eosspark/eos-go/crypto/btcsuite/btcd/btcec"
	"github.com/eosspark/eos-go/crypto/btcsuite/btcutil"
)

type innerK1PrivateKey struct {
	privKey *btcec.PrivateKey
}

func (k *innerK1PrivateKey) publicKey() PublicKey {
	data := k.privKey.PubKey().SerializeCompressed()
	var content [33]byte
	for i := 0; i < 33; {
		content[i] = data[i]
		i += 1
	}
	return PublicKey{Curve: CurveK1, Content: content, inner: &innerK1PublicKey{}}
}

func (k *innerK1PrivateKey) sign(hash []byte) (out Signature, err error) {
	if len(hash) != 32 {
		return out, fmt.Errorf("hash should be 32 bytes")
	}

	compactSig, err := k.privKey.SignCanonical(btcec.S256(), hash)

	if err != nil {
		return out, fmt.Errorf("canonical, %s", err)
	}

	return Signature{Curve: CurveK1, Content: compactSig, innerSignature: &innerK1Signature{}}, nil
}

func (k *innerK1PrivateKey) string() string {
	wif, _ := btcutil.NewWIF(k.privKey, '\x80', false) // no error possible
	return wif.String()
}

func (k *innerK1PrivateKey) Serialize() []byte {
	return k.privKey.Serialize()
}
