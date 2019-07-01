package ecc_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**

$ cleos wallet private_keys -n r1test
password: [[
    "EOS67SCWnz6trqFPCtmxfjYEPSsT9JKRn4zhow8X3VTtgaEzNMULF",
    "5JaKaxySEyjBFGT9K9cYKSFhfojn1RfPcresqRVbmtxnQt1w3qW"
  ],[
    "PUB_R1_6RJ9pXJNe1wk6p2yiJcuJ8QPo7WTudHya9z8vu1VPk44fhBz79",
    "PVT_R1_2o5WfMRU4dTp23pbcbP2yn5MumQzSMy3ayNQ31qi5nUfa2jdWC"
  ],[
    "PUB_R1_7aE3zt3f7cfNuuUwLogDtxSsniQA2uPthATQZ5ErQLuu1nDKFG",
    "PVT_R1_rjKe476v6zXntjC93YAGyqL35NJWshbwcbGRwb27wuKvsRVEa"
  ],[
    "PUB_R1_8KT5dWt33np9V4Nqpdja1GAbkEqVY3pupeYgvCkKTA5FeqePTp",
    "PVT_R1_2FiHVhVjDNjRVAbLg9Cwj1PvVu6Dxn4HKDMFmkyhPZRdAfXwk6"
  ],[
    "PUB_R1_8S4TodyXa9KASMAJgkLbstFYzAWHNjNJPhpHuqqHF9Af8ekV7i",
    "PVT_R1_2sPCnkH6652KFYQZNWuQvgfTTHvqjrhV6pQ8tcVQGqBNsopKZp"
  ]
]

$ echo -n 'banana' | shasum -a 256
b493d48364afe44d11c0165cf470a4164d1e2609911ef998be868d46ade3de4e  -

$ curl --data '["b493d48364afe44d11c0165cf470a4164d1e2609911ef998be868d46ade3de4e","PUB_R1_6RJ9pXJNe1wk6p2yiJcuJ8QPo7WTudHya9z8vu1VPk44fhBz79"]'
http://127.0.0.1:8900/v1/wallet/sign_digest
"SIG_R1_KJmGMknL29w1jTDbkm4wCB5Lr7UXLLWQrfdyurw8dGoTeHggoVbB9wErfUeFhJXwbihuQHK4G4VeaWoNdW7fdScF92Ctx5"

*/

func TestK1PrivateToPublic(t *testing.T) {
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privKey, err := ecc.NewPrivateKey(wif)
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	pubKeyString := pubKey.String()
	assert.Equal(t, "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", pubKeyString)
}

func TestR1PrivateToPublic(t *testing.T) {
	encoded_privKey := "PVT_R1_2o5WfMRU4dTp23pbcbP2yn5MumQzSMy3ayNQ31qi5nUfa2jdWC"
	privKey, err := ecc.NewPrivateKey(encoded_privKey)
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	pubKeyString := pubKey.String()
	assert.Equal(t, "PUB_R1_0000000000000000000000000000000000000000000000", pubKeyString)
}

func TestNewPublicKeyAndSerializeCompress(t *testing.T) {
	// Copied test from eosjs(-.*)?
	key, err := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	require.NoError(t, err)
	assert.Equal(t, "02c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf", hex.EncodeToString(key.Content[:]))
}

func TestNewRandomPrivateKey(t *testing.T) {
	key, err := ecc.NewRandomPrivateKey()
	require.NoError(t, err)
	// taken from eosiojs-ecc:common.test.js:12
	assert.Regexp(t, "^5[HJK].*", key.String())
}

func TestPrivateKeyValidity(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
	}{
		{"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3", true},
		{"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjsm", false},
	}

	for _, test := range tests {
		_, err := ecc.NewPrivateKey(test.in)
		if test.valid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
			assert.Equal(t, "checksum mismatch", err.Error())
		}
	}
}

func TestPublicKeyValidity(t *testing.T) {
	tests := []struct {
		in  string
		err error
	}{
		{"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", nil},
		{"MMM859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM", fmt.Errorf("public key should start with [\"PUB_K1_\" | \"PUB_R1_\"] (or the old \"EOS\")")},
		{"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhTo", fmt.Errorf("checkDecode: invalid checksum")},
	}

	for idx, test := range tests {
		_, err := ecc.NewPublicKey(test.in)
		if test.err == nil {
			assert.NoError(t, err, fmt.Sprintf("test %d with key %q", idx, test.in))
		} else {
			assert.Error(t, err)
			assert.Equal(t, test.err.Error(), err.Error())
		}
	}
}

func TestK1Signature(t *testing.T) {
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privKey, err := ecc.NewPrivateKey(wif)
	require.NoError(t, err)

	cnt := []byte("hi")
	digest := sigDigest([]byte{}, cnt, nil)
	signature, err := privKey.Sign(digest)
	require.NoError(t, err)

	assert.True(t, signature.Verify(digest, privKey.PublicKey()))
}

func TestR1Signature(t *testing.T) {
	encodedPrivKey := "PVT_R1_2o5WfMRU4dTp23pbcbP2yn5MumQzSMy3ayNQ31qi5nUfa2jdWC"
	privKey, err := ecc.NewPrivateKey(encodedPrivKey)
	require.NoError(t, err)

	cnt := []byte("hi")
	digest := sigDigest([]byte{}, cnt, nil)
	_, err = privKey.Sign(digest)
	assert.Error(t, err)
	assert.Equal(t, "R1 not supported", err.Error())
}

func TestNewDeterministicPrivateKey(t *testing.T) {
	a := crypto.Hash256("eosio.token@active")
	g := bytes.NewReader(a.Bytes())
	pri, err := ecc.NewDeterministicPrivateKey(g)
	assert.NoError(t, err)
	assert.Equal(t, "5KNcvkpaba7YDDA9TthNeYybPysBA1aEZJLboRVaYt95NA15nDZ", pri.String())
}

//to do this here because of a import cycle when use eos.SigDigest
func sigDigest(chainID, payload, contextFreeData []byte) []byte {
	h := sha256.New()
	if len(chainID) == 0 {
		_, _ = h.Write(make([]byte, 32, 32))
	} else {
		_, _ = h.Write(chainID)
	}
	_, _ = h.Write(payload)

	if len(contextFreeData) > 0 {
		h2 := sha256.New()
		_, _ = h2.Write(contextFreeData)
		_, _ = h.Write(h2.Sum(nil)) // add the hash of CFD to the payload
	} else {
		_, _ = h.Write(make([]byte, 32, 32))
	}
	return h.Sum(nil)
}
