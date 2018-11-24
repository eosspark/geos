package ecc

import (
	"github.com/eosspark/eos-go/crypto/btcsuite/btcd/btcec"
	"github.com/eosspark/eos-go/crypto/btcsuite/btcutil/base58"
)

type innerK1Signature struct {
}

// verify checks the signature against the pubKey. `hash` is a sha256
// hash of the payload to verify.
func (s *innerK1Signature) verify(content []byte, hash []byte, pubKey PublicKey) bool {
	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)
	if err != nil {
		return false
	}
	key, err := pubKey.Key()
	if err != nil {
		return false
	}
	if recoveredKey.IsEqual(key) {
		return true
	}
	return false
}

func (s *innerK1Signature) publicKey(content []byte, hash []byte) (out PublicKey, err error) {

	recoveredKey, _, err := btcec.RecoverCompact(btcec.S256(), content, hash)

	if err != nil {
		return out, err
	}
	data := recoveredKey.SerializeCompressed()
	var re [33]byte
	for i := 0; i < 33; {
		re[i] = data[i]
		i += 1
	}
	return PublicKey{
		Curve:   CurveK1,
		Content: re,
		inner:   &innerK1PublicKey{},
	}, nil
}

func (s innerK1Signature) string(content []byte) string {
	checksum := Ripemd160checksumHashCurve(content, CurveK1)
	buf := append(content[:], checksum...)
	return "SIG_K1_" + base58.Encode(buf)
}
