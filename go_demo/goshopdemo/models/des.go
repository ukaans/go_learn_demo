package models

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// Des加密
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 使用固定的IV（推荐单独定义，不要直接用key）
	iv := []byte("itying..") // 必须8字节，和key长度一致
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// Des解密
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) == 0 || len(crypted)%8 != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := []byte("itying..") // 和加密时保持一致
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	unpadding := int(origData[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, errors.New("invalid padding size")
	}
	return origData[:(length - unpadding)], nil
}
