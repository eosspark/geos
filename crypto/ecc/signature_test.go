package ecc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureSerialization(t *testing.T) {
	privkey, err := NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	require.NoError(t, err)

	payload := []byte("payload")
	sig, err := privkey.Sign(sigDigest(make([]byte, 32, 32), payload))
	require.NoError(t, err)
	assert.Equal(t, `SIG_K1_K2JjfxmYpoVwCKkohDiQPcepeyetSWMgQPjx3zqagzao5NeQhnW4JQ2qwxd4txU7dR5TdS6PnP75vmMs5qSXzjphpfGnqM`, sig.String()) // not checked after..
	assert.True(t, isCanonical([]byte(sig.Content[:])))
}

func TestSignatureCanonical(t *testing.T) {
	privkey, err := NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	require.NoError(t, err)

	fmt.Println("Start")
	payload := []byte("payload1") // doesn't fail
	sig, err := privkey.Sign(sigDigest(make([]byte, 32, 32), payload))
	fmt.Println("Signed")
	require.NoError(t, err)
	fmt.Println("MAM", sig.String())
	assert.True(t, isCanonical([]byte(sig.Content[:])))
	fmt.Println("End")

	fmt.Println("Start")
	payload = []byte("payload6") // fails
	sig, err = privkey.Sign(sigDigest(make([]byte, 32, 32), payload))
	fmt.Println("Signed")
	require.NoError(t, err)
	fmt.Println("MAM1", sig.String())
	assert.True(t, isCanonical([]byte(sig.Content[:])))
	fmt.Println("End")
}

func TestSignatureMarshalUnmarshal(t *testing.T) {
	fromEOSIOC := "SIG_K1_KYvoFKUyv1xjQYwrAanrRockv9MGJFB55o9SRdobTTi2kMY6PZpEpSz4HxAXNxBTDvbvCcRLY4yA4xBYf542ReXqNvajRi"
	sig, err := NewSignature(fromEOSIOC)
	fmt.Println(err)
	require.NoError(t, err)
	assert.Equal(t, fromEOSIOC, sig.String())
	assert.True(t, isCanonical([]byte(sig.Content[:])))
}
func TestSignatureMarshalUnmarshal_bilc(t *testing.T) {
	fromEOSIOC := "SIG_K1_Jy9G6SgmGSjAbu7n82veUiqV8LFFL6wqr9G26H37dy1WExUj9kYwS17X3ffT5W9M51HkpKF4xQ6MoFCCMxBEHbk64dgbMg"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)
	assert.Equal(t, fromEOSIOC, sig.String())
	assert.True(t, isCanonical([]byte(sig.Content[:])))
}

func isCanonical(compactSig []byte) bool {
	// !(c.data[1] & 0x80)
	// && !(c.data[1] == 0 && !(c.data[2] & 0x80))
	// && !(c.data[33] & 0x80)
	// && !(c.data[33] == 0 && !(c.data[34] & 0x80));

	d := compactSig
	t1 := (d[1] & 0x80) == 0
	t2 := !(d[1] == 0 && ((d[2] & 0x80) == 0))
	t3 := (d[33] & 0x80) == 0
	t4 := !(d[33] == 0 && ((d[34] & 0x80) == 0))
	return t1 && t2 && t3 && t4
}

func TestSignaturePublicKeyExtraction(t *testing.T) {
	//5KX2S6rH1qe4LF14LmRgf5D5C8Kx7QqKynF5cLddC5spe2gtV18
	//EOS5fxEptrpsG2QTjRgi8Gf9EConFDH3jeUc24YFSemcW3bBDhuoW
	fromEOSIOC := "SIG_K1_KfN4vL1cdm51LwGL9jwih6GiMGSe3qTE6EqopdwwyxxJcRJkfoPWpUKjTGrRzspsWhL8qvXQqLjCsBqopN4Z7SPqJJ8UrN"
	sig, err := NewSignature(fromEOSIOC)
	require.NoError(t, err)

	payload, err := hex.DecodeString("20d8af5a0000b32bcc0e37eb0000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500001059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed32320100")
	require.NoError(t, err)
	hashed := crypto.Hash256(payload).Bytes()
	pubKey, err := sig.PublicKey(hashed)
	require.NoError(t, err)
	assert.Equal(t, "EOS5fxEptrpsG2QTjRgi8Gf9EConFDH3jeUc24YFSemcW3bBDhuoW", pubKey.String())

}

func TestEOSIOCSigningComparison(t *testing.T) {
	// try with: ec sign -k 5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3 '{"expiration":"2018-03-21T23:02:32","region":0,"ref_block_num":2156,"ref_block_prefix":1532582828,"packed_bandwidth_words":0,"context_free_cpu_bandwidth":0,"context_free_actions":[],"actions":[],"signatures":[],"context_free_data":[]}'
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3" // corresponds to: EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	chainID, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	payload, err := hex.DecodeString("88e4b25a00006c08ac5b595b000000000000") // without signed transaction bytes
	require.NoError(t, err)

	digest := sigDigest(chainID, payload)

	sig, err := privKey.Sign(digest)
	require.NoError(t, err)

	fromEOSIOC := "SIG_K1_K2WBNtiTY8o4mqFSz7HPnjkiT9JhUYGFa81RrzaXr3aWRF1F8qwVfutJXroqiL35ZiHTcvn8gPWGYJDwnKZTCcbAL56Fxu"
	assert.Equal(t, fromEOSIOC, sig.String())
}

func TestNodeosSignatureComparison(t *testing.T) {
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3" // corresponds to: EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	privKey, err := NewPrivateKey(wif)
	require.NoError(t, err)

	// produce with `cleos create account eosio abourget EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV
	// transaction:
	// chainID + 30d3b35a0000be0194c22fe70000000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed32327c0000000000ea305500000059b1abe93101000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000100000000010000000000ea305500000000a8ed323201000000
	// hashes to:
	digest, _ := hex.DecodeString("a744a49dd60badd5e7073e7287d53e184914242e94ef309d2694e954077dcb27")

	sig, err := privKey.Sign(digest)
	require.NoError(t, err)

	// from that tx:
	fromEOSIOCTx := "SIG_K1_K9JDNpqcgUin9i2PtsV6QbLG8QGzYPN8kqVicJ63CgHBiwq9q27qykaerbNh8kD6baLFWcuKyTmVUwFRF6myjqFQbVsqht"
	assert.Equal(t, fromEOSIOCTx, sig.String())

	// decode
	fmt.Println("From EOSIO sig:", hex.EncodeToString(base58.Decode(fromEOSIOCTx[3:])))
	fmt.Println("From GO sig:", hex.EncodeToString(base58.Decode(sig.String()[3:])))
}

func TestSignatureUnmarshalChecksum(t *testing.T) {
	fromEOSIOC := "SIG_K1_JvnbzEVvC6aZZkuHB75huR9TX8sq6thtLqDTGLg8pmGqRzhAXrMtMJYQqsodtuQ9niBSwS4dEZHdkvWfDsYT9yFQcHvRXf"
	_, err := NewSignature(fromEOSIOC)
	require.Equal(t, "signature checksum failed, found 490cc16d expected 141e1d3a", err.Error())
}

func sigDigest(chainID, payload []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(chainID)
	_, _ = h.Write(payload)
	return h.Sum(nil)
}

func TestSignatureMarshal(t *testing.T) {
	fromEOSIOC := "SIG_K1_JvnbzEVvC6aZZkuHB75huR9TX8sq6thtLqDTGLg8pmGqRzhAXrMtMJYQqsodtuQ9niBSwS4dEZHdkvWfDsYT9yvuwgM4Bi"
	sig, err := NewSignature(fromEOSIOC)
	fmt.Println(sig, err)
	require.NoError(t, err)
	assert.Equal(t, fromEOSIOC, sig.String())
	assert.True(t, isCanonical([]byte(sig.Content[:])))

}
