package ecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/eosspark/eos-go/exception/try"
	"strings"

	"github.com/eosspark/eos-go/crypto/btcsuite/btcd/btcec"
	"github.com/eosspark/eos-go/crypto/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"reflect"
)

const PublicKeyPrefix = "PUB_"
const PublicKeyK1Prefix = "PUB_K1_"
const PublicKeyR1Prefix = "PUB_R1_"
const PublicKeyPrefixCompat = "EOS"

type innerPublicKey interface {
	key(content []byte) (*btcec.PublicKey, error)
	string(content []byte, curveID CurveID) string
}

type PublicKey struct {
	Curve   CurveID
	Content [33]byte

	inner innerPublicKey
}

func NewPublicKeyFromData(data []byte) (out PublicKey, err error) {
	if len(data) != 34 {
		return out, fmt.Errorf("public key data must have a length of 33 ")
	}
	var content [33]byte
	for i := 0; i < 33; {
		content[i] = data[i+1]
		i += 1
	}
	out = PublicKey{
		Curve:   CurveID(data[0]), // 1 byte
		Content: content,          // 33 bytes
	}

	switch out.Curve {
	case CurveK1:
		out.inner = &innerK1PublicKey{}
	case CurveR1:
		out.inner = &innerR1PublicKey{}
	default:
		return out, fmt.Errorf("unsupported curve prefix %q", out.Curve)
	}

	return out, nil
}

func MustNewPublicKeyFromData(data []byte) PublicKey {
	key, err := NewPublicKeyFromData(data)
	Throw(err)
	return key
}

func NewPublicKey(pubKey string) (out PublicKey, err error) {
	if len(pubKey) < 8 {
		return out, fmt.Errorf("invalid format")
	}

	var decodedPubKey []byte
	var curveID CurveID
	var inner innerPublicKey

	if strings.HasPrefix(pubKey, PublicKeyR1Prefix) {
		pubKeyMaterial := pubKey[len(PublicKeyR1Prefix):] // strip "PUB_R1_"
		decodedPubKey = base58.Decode(pubKeyMaterial)
		inner = &innerR1PublicKey{}
	} else if strings.HasPrefix(pubKey, PublicKeyK1Prefix) {
		pubKeyMaterial := pubKey[len(PublicKeyK1Prefix):] // strip "PUB_K1_"
		curveID = CurveK1
		decodedPubKey, err = checkDecode(pubKeyMaterial, curveID)
		if err != nil {
			return out, fmt.Errorf("checkDecode: %s", err)
		}
		inner = &innerK1PublicKey{}
	} else if strings.HasPrefix(pubKey, PublicKeyPrefixCompat) { // "EOS"
		pubKeyMaterial := pubKey[len(PublicKeyPrefixCompat):] // strip "EOS"
		curveID = CurveK1
		decodedPubKey, err = checkDecode(pubKeyMaterial, curveID)
		if err != nil {
			return out, fmt.Errorf("checkDecode: %s", err)
		}
		inner = &innerK1PublicKey{}
	} else {
		return out, fmt.Errorf("public key should start with [%q | %q] (or the old %q)", PublicKeyK1Prefix, PublicKeyR1Prefix, PublicKeyPrefixCompat)
	}

	var content [33]byte
	for i := 0; i < 33; {
		content[i] = decodedPubKey[i]
		i += 1
	}
	return PublicKey{Curve: curveID, Content: content, inner: inner}, nil
}

func MustNewPublicKey(pubKey string) PublicKey {
	key, err := NewPublicKey(pubKey)
	Throw(err)
	return key
}
func NewPublicKeyNil() *PublicKey {
	return &PublicKey{Curve: CurveK1, Content: [33]byte{}, inner: &innerK1PublicKey{}}

}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
func checkDecode(input string, curve CurveID) (result []byte, err error) {
	decoded := base58.Decode(input)
	if len(decoded) < 5 {
		return nil, fmt.Errorf("invalid format")
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	///// WARN: ok the ripemd160checksum should include the prefix in CERTAIN situations,
	// like when we imported the PubKey without a prefix ?! tied to the string representation
	// or something ? weird.. checksum shouldn't change based on the string reprsentation.
	if bytes.Compare(ripemd160checksum(decoded[:len(decoded)-4], curve), cksum[:]) != 0 {
		return nil, fmt.Errorf("invalid checksum")
	}
	// perhaps bitcoin has a leading net ID / version, but EOS doesn't
	payload := decoded[:len(decoded)-4]
	result = append(result, payload...)
	return
}

func ripemd160checksum(in []byte, curve CurveID) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path

	// if curve != CurveK1 {
	// 	_, _ = h.Write([]byte(curve.String())) // conditionally ?
	// }
	sum := h.Sum(nil)
	return sum[:4]
}

func Ripemd160checksumHashCurve(in []byte, curve CurveID) []byte {
	h := ripemd160.New()
	_, _ = h.Write(in) // this implementation has no error path

	// FIXME: this seems to be only rolled out to the `SIG_` things..
	// proper support for importing `EOS` keys isn't rolled out into `dawn4`.
	_, _ = h.Write([]byte(curve.String())) // conditionally ?
	sum := h.Sum(nil)
	return sum[:4]
}

func (p PublicKey) Key() (*btcec.PublicKey, error) {
	return p.inner.key(p.Content[:])
}

func (p PublicKey) String() string {

	hash := ripemd160checksum(p.Content[:], p.Curve)

	rawKey := make([]byte, 37)
	copy(rawKey, p.Content[:])
	copy(rawKey[33:], hash[:4])

	return PublicKeyPrefixCompat + base58.Encode(rawKey)
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	s := p.String()
	return json.Marshal(s)
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	newKey, err := NewPublicKey(s)
	if err != nil {
		return err
	}

	*p = newKey

	return nil
}

func (p PublicKey) Valid() bool {
	switch p.Curve {
	case CurveK1:
		_, err := btcec.ParsePubKey(p.Content[:], btcec.S256())
		if err != nil {
			return false
		}
	case CurveR1:
		_, err := btcec.ParsePubKey(p.Content[:], btcec.S256R1())
		if err != nil {
			return false
		}
	default:
		return false
	}
	return true
}

func (p PublicKey) Compare(pub PublicKey) bool {
	if p.Curve != pub.Curve {
		return false
	}
	if p.Content != pub.Content {
		return false
	}
	return true
}

func (p PublicKey) GetKey() []byte {
	sl, _ := p.MarshalJSON()
	return sl
}

var TypePubKey = reflect.TypeOf(PublicKey{})

func ComparePubKey(first interface{}, second interface{}) int {
	pub1 := first.(PublicKey)
	pub2 := second.(PublicKey)

	if comp := bytes.Compare(pub1.Content[:], pub2.Content[:]); comp != 0 {
		return comp
	}

	if pub1.Curve > pub2.Curve {
		return 1
	} else if pub1.Curve < pub2.Curve {
		return -1
	} else {
		return 0
	}
}
