// Copyright 2021 dudaodong@gmail.com. All rights reserved.
// Use of this source code is governed by MIT license

// Package cryptor implements some util functions to encrypt and decrypt.
// Note:
// 1. for aes crypt function, the `key` param length should be 16, 24 or 32. if not, will err = errors.New.
package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
)

// AesEcbEncrypt encrypt data with key use AES ECB algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/jT5irszHx-j
func AesEcbEncrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}

	length := (len(data) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)

	copy(plain, data)

	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}

	encrypted := make([]byte, len(plain))
	cipher, err := aes.NewCipher(generateAesKey(key, size))
	if err != nil {
		return
	}
	for bs, be := 0, cipher.BlockSize(); bs <= len(data); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted, err
}

// AesEcbDecrypt decrypt data with key use AES ECB algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/jT5irszHx-j
func AesEcbDecrypt(encrypted, key []byte) (r []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}
	cipher, err := aes.NewCipher(generateAesKey(key, size))
	if err != nil {
		return
	}
	decrypted := make([]byte, len(encrypted))

	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim], nil
}

// AesCbcEncrypt encrypt data with key use AES CBC algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/IOq_g8_lKZD
func AesCbcEncrypt(data, key []byte) (r []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	data = pkcs7Padding(data, block.BlockSize())

	encrypted := make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted[aes.BlockSize:], data)

	return encrypted, nil
}

// AesCbcDecrypt decrypt data with key use AES CBC algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/IOq_g8_lKZD
func AesCbcDecrypt(encrypted, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encrypted, encrypted)

	decrypted := pkcs7UnPadding(encrypted)
	return decrypted, nil
}

// AesCtrCrypt encrypt data with key use AES CTR algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/SpaZO0-5Nsp
func AesCtrCrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}

	block, _ := aes.NewCipher(key)

	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	dst := make([]byte, len(data))
	stream.XORKeyStream(dst, data)

	return dst, nil
}

// AesCfbEncrypt encrypt data with key use AES CFB algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/tfkF10B13kH
func AesCfbEncrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	encrypted := make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], data)

	return encrypted, nil
}

// AesCfbDecrypt decrypt data with key use AES CFB algorithm
// len(encrypted) should be great than 16, len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/tfkF10B13kH
func AesCfbDecrypt(encrypted, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}

	if len(encrypted) < aes.BlockSize {
		err = errors.New("encrypted data is too short")
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(encrypted, encrypted)

	return encrypted, err
}

// AesOfbEncrypt encrypt data with key use AES OFB algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/VtHxtkUj-3F
func AesOfbEncrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	data = pkcs7Padding(data, aes.BlockSize)
	encrypted := make([]byte, aes.BlockSize+len(data))
	iv := encrypted[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], data)

	return encrypted, err
}

// AesOfbDecrypt decrypt data with key use AES OFB algorithm
// len(key) should be 16, 24 or 32.
// Play: https://go.dev/play/p/VtHxtkUj-3F
func AesOfbDecrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		err = errors.New("key length shoud be 16 or 24 or 32")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	if len(data)%aes.BlockSize != 0 {
		err = errors.New("data empty")
		return
	}

	decrypted := make([]byte, len(data))
	mode := cipher.NewOFB(block, iv)
	mode.XORKeyStream(decrypted, data)

	decrypted = pkcs7UnPadding(decrypted)

	return decrypted, nil
}

// DesEcbEncrypt encrypt data with key use DES ECB algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/8qivmPeZy4P
func DesEcbEncrypt(data, key []byte) (bts []byte, err error) {
	length := (len(data) + des.BlockSize) / des.BlockSize
	plain := make([]byte, length*des.BlockSize)
	copy(plain, data)

	pad := byte(len(plain) - len(data))
	for i := len(data); i < len(plain); i++ {
		plain[i] = pad
	}

	encrypted := make([]byte, len(plain))
	cipher, err := des.NewCipher(generateDesKey(key))
	if err != nil {
		return
	}
	for bs, be := 0, cipher.BlockSize(); bs <= len(data); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted, err
}

// DesEcbDecrypt decrypt data with key use DES ECB algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/8qivmPeZy4P
func DesEcbDecrypt(encrypted, key []byte) (bts []byte, err error) {
	cipher, err := des.NewCipher(generateDesKey(key))
	if err != nil {
		return
	}
	decrypted := make([]byte, len(encrypted))

	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim], err
}

func errorTest() (r int, err error) {
	return 1, errors.New("JustError")
}

func ErrorTest() (err error) {
	if a, err := errorTest(); err != nil {
		return err
	} else {
		fmt.Println(a)
		return nil
	}
}

// DesCbcEncrypt encrypt data with key use DES CBC algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/4cC4QvWfe3_1
func DesCbcEncrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 8 {
		err = errors.New("key length shoud be 8")
		return
	}

	block, _ := des.NewCipher(key)
	data = pkcs7Padding(data, block.BlockSize())

	encrypted := make([]byte, des.BlockSize+len(data))
	iv := encrypted[:des.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted[des.BlockSize:], data)

	return encrypted, nil
}

// DesCbcDecrypt decrypt data with key use DES CBC algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/4cC4QvWfe3_1
func DesCbcDecrypt(encrypted, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 8 {
		err = errors.New("key length shoud be 8")
		return
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return
	}
	iv := encrypted[:des.BlockSize]
	encrypted = encrypted[des.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encrypted, encrypted)

	decrypted := pkcs7UnPadding(encrypted)
	return decrypted, err
}

// DesCtrCrypt encrypt data with key use DES CTR algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/9-T6OjKpcdw
func DesCtrCrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 8 {
		err = errors.New("key length shoud be 8")
		return
	}

	block, err := des.NewCipher(key)

	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	dst := make([]byte, len(data))
	stream.XORKeyStream(dst, data)

	return dst, err
}

// DesCfbEncrypt encrypt data with key use DES CFB algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/y-eNxcFBlxL
func DesCfbEncrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 8 {
		err = errors.New("key length shoud be 8")
		return
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return
	}

	encrypted := make([]byte, des.BlockSize+len(data))
	iv := encrypted[:des.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[des.BlockSize:], data)

	return encrypted, err
}

// DesCfbDecrypt decrypt data with key use DES CFB algorithm
// len(encrypted) should be great than 16, len(key) should be 8.
// Play: https://go.dev/play/p/y-eNxcFBlxL
func DesCfbDecrypt(encrypted, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 8 {
		err = errors.New("key length shoud be 8")
		return
	}

	block, err := des.NewCipher(key)
	if len(encrypted) < des.BlockSize {
		err = errors.New("encrypted data is too short")
		return
	}
	iv := encrypted[:des.BlockSize]
	encrypted = encrypted[des.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)

	return encrypted, err
}

// DesOfbEncrypt encrypt data with key use DES OFB algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/74KmNadjN1J
func DesOfbEncrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 8 {
		err = errors.New("key length shoud be 8")
		return
	}

	block, err := des.NewCipher(key)
	if err != nil {

		return
	}
	data = pkcs7Padding(data, des.BlockSize)
	encrypted := make([]byte, des.BlockSize+len(data))
	iv := encrypted[:des.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(encrypted[des.BlockSize:], data)

	return encrypted, err
}

// DesOfbDecrypt decrypt data with key use DES OFB algorithm
// len(key) should be 8.
// Play: https://go.dev/play/p/74KmNadjN1J
func DesOfbDecrypt(data, key []byte) (bts []byte, err error) {
	size := len(key)
	if size != 8 {
		err = errors.New("key length shoud be 8")
		return
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return
	}

	iv := data[:des.BlockSize]
	data = data[des.BlockSize:]
	if len(data)%des.BlockSize != 0 {
		err = errors.New("data len is wrong")
		return
	}

	decrypted := make([]byte, len(data))
	mode := cipher.NewOFB(block, iv)
	mode.XORKeyStream(decrypted, data)

	decrypted = pkcs7UnPadding(decrypted)

	return decrypted, err
}

// GenerateRsaKey create rsa private and public pemo file.
// Play: https://go.dev/play/p/zutRHrDqs0X
func GenerateRsaKey(keySize int, priKeyFile, pubKeyFile string) error {
	// private key
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}

	derText := x509.MarshalPKCS1PrivateKey(privateKey)

	block := pem.Block{
		Type:  "rsa private key",
		Bytes: derText,
	}

	file, err := os.Create(priKeyFile)
	if err != nil {
		return err
	}
	err = pem.Encode(file, &block)
	if err != nil {
		return err
	}

	file.Close()

	// public key
	publicKey := privateKey.PublicKey

	derpText, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}

	block = pem.Block{
		Type:  "rsa public key",
		Bytes: derpText,
	}

	file, err = os.Create(pubKeyFile)
	if err != nil {
		return err
	}

	err = pem.Encode(file, &block)
	if err != nil {
		return err
	}

	file.Close()

	return nil
}

// RsaEncrypt encrypt data with ras algorithm.
// Play: https://go.dev/play/p/rDqTT01SPkZ
func RsaEncrypt(data []byte, pubKeyFileName string) (bts []byte, err error) {
	file, err := os.Open(pubKeyFileName)
	if err != nil {

		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	defer file.Close()
	buf := make([]byte, fileInfo.Size())

	_, err = file.Read(buf)
	if err != nil {
		return
	}

	block, _ := pem.Decode(buf)

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pubKey := pubInterface.(*rsa.PublicKey)

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return
	}
	return cipherText, err
}

// RsaDecrypt decrypt data with ras algorithm.
// Play: https://go.dev/play/p/rDqTT01SPkZ
func RsaDecrypt(data []byte, privateKeyFileName string) (bts []byte, err error) {
	file, err := os.Open(privateKeyFileName)
	if err != nil {
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	buf := make([]byte, fileInfo.Size())
	defer file.Close()

	_, err = file.Read(buf)
	if err != nil {
		return
	}

	block, _ := pem.Decode(buf)

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, data)
	if err != nil {
		return
	}
	return plainText, err
}

// GenerateRsaKeyPair create rsa private and public key.
// Play: https://go.dev/play/p/sSVmkfENKMz
func GenerateRsaKeyPair(keySize int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, keySize)
	return privateKey, &privateKey.PublicKey
}

// RsaEncryptOAEP encrypts the given data with RSA-OAEP.
// Play: https://go.dev/play/p/sSVmkfENKMz
func RsaEncryptOAEP(data []byte, label []byte, key rsa.PublicKey) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key, data, label)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

// RsaDecryptOAEP decrypts the data with RSA-OAEP.
// Play: https://go.dev/play/p/sSVmkfENKMz
func RsaDecryptOAEP(ciphertext []byte, label []byte, key rsa.PrivateKey) ([]byte, error) {
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &key, ciphertext, label)
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}
