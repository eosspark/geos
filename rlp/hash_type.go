package rlp

import (
	"bytes"
)

const hashingChecktimeBlockSize uint32 = 10 * 1024

type ShaType interface {
	Write(p []byte) (nn int, err error)
	Sum(b []byte) []byte
}

func encode(sha ShaType, data []byte, l uint32) []byte {
	var pos uint32 = 0
	for l > hashingChecktimeBlockSize {
		_, _ = sha.Write(data[pos : pos+hashingChecktimeBlockSize])
		pos += hashingChecktimeBlockSize
		l -= hashingChecktimeBlockSize
		//checktime()
	}
	sha.Write(data[pos:l])
	return sha.Sum(nil)
}

func AssertSha1(data []byte, datalen uint32, hashVal *Sha1) {
	h := NewSha1()
	hash := encode(h, data, datalen)
	if bytes.Compare(hash, hashVal.Bytes()) != 0 {
		panic("hash mismatch")
	}
}

func AssertSha256(data []byte, datalen uint32, hashVal *Sha256) {
	h := NewSha256()
	hash := encode(h, data, datalen)
	if bytes.Compare(hash, hashVal.Bytes()) != 0 {
		panic("hash mismatch")
	}
}

func AssertSha512(data []byte, datalen uint32, hashVal *Sha512) {
	h := NewSha512()
	hash := encode(h, data, datalen)
	if bytes.Compare(hash, hashVal.Bytes()) != 0 {
		panic("hash mismatch")
	}
}

func AssertRipemd160(data []byte, datalen uint32, hashVal *Ripemd160) {
	h := NewRipemd160()
	hash := encode(h, data, datalen)
	if bytes.Compare(hash, hashVal.Bytes()) != 0 {
		panic("hash mismatch")
	}
}

func Sha1Hash(data []byte, datalen uint32, hash_val *Sha1) {
	h := NewSha1()
	hash := encode(h, data, datalen)
	*hash_val = *NewSha1Byte(hash)
}

func Sha256Hash(data []byte, datalen uint32, hash_val *Sha256) {
	h := NewSha256()
	hash := encode(h, data, datalen)
	*hash_val = *NewSha256Byte(hash)
}
func Sha512Hash(data []byte, datalen uint32, hash_val *Sha512) {
	h := NewSha512()
	hash := encode(h, data, datalen)
	*hash_val = *NewSha512Byte(hash)
}
func Ripemd160Hash(data []byte, datalen uint32, hash_val *Ripemd160) {
	h := NewRipemd160()
	hash := encode(h, data, datalen)
	*hash_val = *NewRipemd160Byte(hash)
}
