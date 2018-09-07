package walletPlugin

import (
	"crypto/aes"
	"crypto/cipher"
)

type AesEncrypt struct{}

func (a *AesEncrypt) getKey(str string) []byte {
	strKey := str
	keyLen := len(strKey)
	if keyLen < 16 {
		panic("res key lenth is < 16")
	}
	arrKey := []byte(strKey)
	if keyLen > 32 {
		return arrKey[:32]
	}
	if keyLen >= 24 {
		return arrKey[0:24]
	}
	return arrKey[:16]
}

func (a *AesEncrypt) Encrypt(keystr string, strMsg string) ([]byte, error) {
	key := a.getKey(keystr)
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMsg))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMsg))
	return encrypted, nil
}

func (a *AesEncrypt) Decrypt(keystr string, src []byte) (strDesc []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	key := a.getKey(keystr)

	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var aesBlockDerypter cipher.Block
	aesBlockDerypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDerypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return decrypted, nil

}

// func Cipher_test() {
//   aesEnc := AesEncrypt{}
//   arrEncrypt, err := aesEnc.Encrypt("PW5JqYnp8DtqzLYr8wFxZuiJWTfmygPDvKaF1U45hkuL5yo68mZJ6", "5KTiH48Xj1onYwPjRjHNPW6MNwJaVeqDk5f3u3cQ3jUasc3u75f")
//   if err != nil {
//     fmt.Println(arrEncrypt)
//     return
//   }
//   fmt.Println(arrEncrypt)
//   strMsg, err := aesEnc.Decrypt("PW5JqYnp8DtqzLYr8wFxZuiJWTfmygPDvKaF1U45hkuL5yo68mZJ6", arrEncrypt)
//   if err != nil {
//     fmt.Println(arrEncrypt)
//     return
//   }
//   fmt.Println(string(strMsg))
// }
