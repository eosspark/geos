package exec

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"strings"
)

// void assert_recover_key( const fc::sha256& digest,
//                   array_ptr<char> sig, size_t siglen,
//                   array_ptr<char> pub, size_t publen ) {
//    fc::crypto::signature s;
//    fc::crypto::public_key p;
//    datastream<const char*> ds( sig, siglen );
//    datastream<const char*> pubds( pub, publen );

//    fc::raw::unpack(ds, s);
//    fc::raw::unpack(pubds, p);

//    auto check = fc::crypto::public_key( s, digest, false );
//    EOS_ASSERT( check == p, crypto_api_exception, "Error expected key different than recovered key" );
// }
func assertRecoverKey(w *WasmInterface, digest int,
	sig int, siglen int,
	pub int, publen int) {
	fmt.Println("assert_recover_key")

	digBytes := getSha256(w, digest)
	sigBytes := getMemory(w, sig, siglen)
	pubBytes := getMemory(w, pub, publen)

	fmt.Println("d:", hex.EncodeToString(digBytes), " s:", hex.EncodeToString(sigBytes), " p:", hex.EncodeToString(pubBytes))

	var s ecc.Signature
	var p ecc.PublicKey
	//var d []byte

	//rlp.DecodeBytes(digBytes, &d)
	d := digBytes
	rlp.DecodeBytes(sigBytes, &s)
	rlp.DecodeBytes(pubBytes, &p)

	check, err := s.PublicKey(d)
	if err != nil {
		return
		//assert
	}

	if strings.Compare(check.String(), p.String()) != 0 {
		println("Error expected key different than recovered key")
		//assert
	}

}

// int recover_key( const fc::sha256& digest,
//                   array_ptr<char> sig, size_t siglen,
//                   array_ptr<char> pub, size_t publen ) {
//    fc::crypto::signature s;
//    datastream<const char*> ds( sig, siglen );
//    datastream<char*> pubds( pub, publen );

//    fc::raw::unpack(ds, s);
//    fc::raw::pack( pubds, fc::crypto::public_key( s, digest, false ) );
//    return pubds.tellp();
// }
func recoverKey(w *WasmInterface, digest int,
	sig int, siglen int,
	pub int, publen int) int {
	fmt.Println("recover_key")

	digBytes := getSha256(w, digest)
	sigBytes := getMemory(w, sig, siglen)

	var s ecc.Signature
	rlp.DecodeBytes(sigBytes, &s)
	check, _ := s.PublicKey(digBytes)

	p, err := rlp.EncodeToBytes(check)
	if err != nil {
		return -1
	}

	l := len(p)
	if l > publen {
		l = publen
	}
	setMemory(w, pub, p, 0, l)

	return l
}

// template<class Encoder> auto encode(char* data, size_t datalen) {
//    Encoder e;
//    const size_t bs = eosio::chain::config::hashing_checktime_block_size;
//    while ( datalen > bs ) {
//       e.write( data, bs );
//       data += bs;
//       datalen -= bs;
//       context.trx_context.checktime();
//    }
//    e.write( data, datalen );
//    return e.result();
// }

// void assert_sha256(array_ptr<char> data, size_t datalen, const fc::sha256& hash_val) {
//    auto result = encode<fc::sha256::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
var (
	hashingChecktimeBlockSize uint = 10 * 1024 //move to config
)

type shaInterface interface {
	Write(p []byte) (nn int, err error)
	Sum(b []byte) []byte
}

func encode(w *WasmInterface, s shaInterface, data []byte, dataLen int) []byte {

	bs := int(hashingChecktimeBlockSize)

	i := 0
	l := dataLen

	for i = 0; l > bs; i += bs {
		s.Write(data[i : i+bs])
		l -= bs
		w.context.CheckTime()
	}

	s.Write(data[i : i+l])

	return s.Sum(nil)

}

func assertSha256(w *WasmInterface, data int, datalen int, hash_val int) {
	fmt.Println("assert_sha256")

	dataBytes := getMemory(w, data, datalen)
	if dataBytes == nil {
		return
	}
	//var s rlp.Sha256
	s := rlp.NewSha256()
	hashEncode := encode(w, s, dataBytes, datalen)
	hash := getSha256(w, hash_val)

	if bytes.Compare(hashEncode, hash) != 0 {
		fmt.Println("sha256 hash mismatch")
		//assert
	}
}

// void assert_sha1(array_ptr<char> data, size_t datalen, const fc::sha1& hash_val) {
//    auto result = encode<fc::sha1::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
func assertSha1(w *WasmInterface, data int, dataLen int, hash_val int) {
	fmt.Println("assert_sha1")

	dataBytes := getMemory(w, data, dataLen)
	if dataBytes == nil {
		return
	}

	//var s crypto.Sha1
	//s := sha1.New()
	s := crypto.NewSha1()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := getSha1(w, hash_val)

	if bytes.Compare(hashEncode, hash) != 0 {
		fmt.Println("sha1 hash mismatch")
	}
}

// void assert_sha512(array_ptr<char> data, size_t datalen, const fc::sha512& hash_val) {
//    auto result = encode<fc::sha512::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
func assertSha512(w *WasmInterface, data int, dataLen int, hash_val int) {
	fmt.Println("assert_sha512")

	dataBytes := getMemory(w, data, dataLen)
	if dataBytes == nil {
		return
	}

	s := rlp.NewSha512()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := getSha512(w, hash_val)

	if bytes.Compare(hashEncode, hash) != 0 {
		fmt.Println("sha512 hash mismatch")
		//assert
	}

}

// void assert_ripemd160(array_ptr<char> data, size_t datalen, const fc::ripemd160& hash_val) {
//    auto result = encode<fc::ripemd160::encoder>( data, datalen );
//    EOS_ASSERT( result == hash_val, crypto_api_exception, "hash mismatch" );
// }
func assertRipemd160(w *WasmInterface, data int, dataLen int, hash_val int) {
	fmt.Println("assert_ripemd160")

	dataBytes := getMemory(w, data, dataLen)
	if dataBytes == nil {
		return
	}

	s := rlp.NewRipemd160()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := getRipemd160(w, hash_val)

	if bytes.Compare(hashEncode, hash) != 0 {
		fmt.Println("ripemd160 hash mismatch")
		//assert
	}
}

// void sha1(array_ptr<char> data, size_t datalen, fc::sha1& hash_val) {
//    hash_val = encode<fc::sha1::encoder>( data, datalen );
// }
func sha1(w *WasmInterface, data int, dataLen int, hash_val int) {
	fmt.Println("sha1")

	dataBytes := getMemory(w, data, dataLen)
	if dataBytes == nil {
		return
	}

	s := rlp.NewSha1()
	hashEncode := encode(w, s, dataBytes, dataLen)
	setSha1(w, hash_val, hashEncode)
}

// void sha256(array_ptr<char> data, size_t datalen, fc::sha256& hash_val) {
//    hash_val = encode<fc::sha256::encoder>( data, datalen );
// }
func sha256(w *WasmInterface, data int, dataLen int, hash_val int) {
	fmt.Println("sha256")

	dataBytes := getMemory(w, data, dataLen)
	if dataBytes == nil {
		return
	}

	s := rlp.NewSha256()

	hashEncode := encode(w, s, dataBytes, dataLen)
	setSha256(w, hash_val, hashEncode)
}

// void sha512(array_ptr<char> data, size_t datalen, fc::sha512& hash_val) {
//    hash_val = encode<fc::sha512::encoder>( data, datalen );
// }
func sha512(w *WasmInterface, data int, dataLen int, hash_val int) {
	fmt.Println("sha512")

	dataBytes := getMemory(w, data, dataLen)
	if dataBytes == nil {
		return
	}

	s := rlp.NewSha512()

	hashEncode := encode(w, s, dataBytes, dataLen)
	setSha512(w, hash_val, hashEncode)
}

// void ripemd160(array_ptr<char> data, size_t datalen, fc::ripemd160& hash_val) {
//    hash_val = encode<fc::ripemd160::encoder>( data, datalen );
// }
func ripemd160(w *WasmInterface, data int, dataLen int, hash_val int) {
	fmt.Println("ripemd160")

	dataBytes := getMemory(w, data, dataLen)
	if dataBytes == nil {
		return
	}

	s := rlp.NewRipemd160()
	hashEncode := encode(w, s, dataBytes, dataLen)
	setRipemd160(w, hash_val, hashEncode)
}
