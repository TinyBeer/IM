package datasafe

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

var key []byte = []byte("asdfghjklqwertyu")

// 加密
func EncryptoAES(src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// 解密
func DecryptoAES(src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}

// 填充数据
func padding(src []byte, blockSize int) []byte {
	// 计算需要填充的数量
	padNum := blockSize - len(src)%blockSize

	// 准备填充的数据  使用填充的需要填充数填充，方便解码
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)

	// 返回填充结果
	return append(src, pad...)
}

// 去掉填充的数据
func unpadding(src []byte) []byte {
	n := len(src)
	unpadNum := int(src[n-1])

	return src[:n-unpadNum]
}
