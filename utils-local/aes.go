package utils_local

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"bytes"
	"errors"
	"strings"
	"fmt"
)

func AesEncrypt(origData string, key string, slash string) string {
	keyByte := []byte(key)
	origDataByte := []byte(origData)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return ""
	}
	blockSize := block.BlockSize()
	origDataByte = PKCS5Padding(origDataByte, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])
	crypted := make([]byte, len(origDataByte))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origDataByte)
	cipherText := base64.StdEncoding.EncodeToString(crypted)
	cipherText = strings.Replace(cipherText, "/", slash, -1)
	return cipherText
}

func AesDecrypt(cipherText string, key string, slash string) (s string, e error) {
	if cipherText == "" {
		return "", errors.New("empty cipher text")
	}
	defer func() {
		if r := recover(); r != nil{
			e = errors.New(fmt.Sprint(r))
		}
	}()
	cipherText = strings.Replace(cipherText, slash, "/", -1)
	keyByte := []byte(key)
	cryptedByte, _ := base64.StdEncoding.DecodeString(cipherText)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyByte[:blockSize])
	origData := make([]byte, len(cryptedByte))
	// origData := crypted
	blockMode.CryptBlocks(origData, cryptedByte)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

