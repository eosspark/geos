package ecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/eosspark/eos-go/exception/try"
	"strings"

	"github.com/eosspark/eos-go/crypto/btcsuite/btcutil/base58"
)

type innerSignature interface {
	verify(content []byte, hash []byte, pubKey PublicKey) bool
	publicKey(content []byte, hash []byte) (out PublicKey, err error)
	string(content []byte) string
}

// Signature represents a signature for some hash
type Signature struct {
	Curve   CurveID
	Content []byte // the Compact signature as bytes

	innerSignature innerSignature
}

func (s Signature) Pack() ([]byte, error) {
	re := make([]byte, 0)

	if len(s.Content) == 0 {
		s.Curve = CurveK1
		s.Content = make([]byte, 65)
	}
	if len(s.Content) != 65 {
		return nil, fmt.Errorf("signature should be 65 bytes, was %d", len(s.Content))
	}
	re = append(re, byte(s.Curve))
	re = append(re, s.Content...)
	return re, nil

}
func (s *Signature) Unpack(in []byte) (l int, err error) {
	if len(in) < 66 {
		return 0, fmt.Errorf("signature required [%d] bytes, remaining [%d]", 66, len(in))
	}
	sigContent := make([]byte, 66)
	copy(sigContent, in[:66])
	*s, err = NewSignatureFromData(sigContent)
	if err != nil {
		return 0, fmt.Errorf("new signature: %s", err)
	}
	return 66, nil

}

func (s Signature) Verify(hash []byte, pubKey PublicKey) bool {
	return s.innerSignature.verify(s.Content, hash, pubKey)
}

func (s Signature) PublicKey(hash []byte) (out PublicKey, err error) {
	return s.innerSignature.publicKey(s.Content, hash)
}

func (s Signature) String() string {
	return s.innerSignature.string(s.Content)
}

func NewSignatureFromData(data []byte) (Signature, error) {
	if len(data) != 66 {
		return Signature{}, fmt.Errorf("data length of a signature should be 66, reveived %d", len(data))
	}

	signature := Signature{
		Curve:   CurveID(data[0]), // 1 byte
		Content: data[1:],         // 65 bytes
	}

	switch signature.Curve {
	case CurveK1:
		signature.innerSignature = &innerK1Signature{}
	case CurveR1:
		signature.innerSignature = &innerR1Signature{}
	default:
		return Signature{}, fmt.Errorf("invalid curve  %q", signature.Curve)
	}
	return signature, nil
}

func MustNewSignatureFromData(data []byte) Signature {
	sig, err := NewSignatureFromData(data)
	Throw(err)
	return sig
}

func NewSignature(fromText string) (Signature, error) {
	if !strings.HasPrefix(fromText, "SIG_") {
		return Signature{}, fmt.Errorf("signature should start with SIG_")
	}
	if len(fromText) < 8 {
		return Signature{}, fmt.Errorf("invalid signature length")
	}

	fromText = fromText[4:] // remove the `SIG_` prefix

	var curvePrefix = fromText[:3]
	switch curvePrefix {
	case "K1_":

		fromText = fromText[3:] // strip curve ID

		sigbytes := base58.Decode(fromText)

		content := sigbytes[:len(sigbytes)-4]
		checksum := sigbytes[len(sigbytes)-4:]
		verifyChecksum := Ripemd160checksumHashCurve(content, CurveK1)
		if !bytes.Equal(verifyChecksum, checksum) {
			return Signature{}, fmt.Errorf("signature checksum failed, found %x expected %x", verifyChecksum, checksum)
		}

		return Signature{Curve: CurveK1, Content: content, innerSignature: &innerK1Signature{}}, nil

	case "R1_":

		fromText = fromText[3:] // strip R1_
		content := base58.Decode(fromText)
		//todo: stuff here

		return Signature{Curve: CurveR1, Content: content, innerSignature: &innerR1Signature{}}, nil

	default:
		return Signature{}, fmt.Errorf("invalid curve prefix %q", curvePrefix)
	}
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Signature) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		s = NewSigNil()
		return nil
	}
	var sig string
	err = json.Unmarshal(data, &sig)
	if err != nil {
		return
	}

	*s, err = NewSignature(sig)

	return
}

func NewSigNil() *Signature {
	return &Signature{
		Curve:          CurveK1,
		Content:        make([]byte, 65, 65),
		innerSignature: &innerK1Signature{},
	}
}
