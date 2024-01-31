package openssl

import (
	"bytes"
	"crypto/cipher"
)

// CBCEncrypt 使用CBC模式加密算法对给定的数据进行加密
// 参数：
//   - block：使用的加密块
//   - src：要加密的数据
//   - iv：初始化向量（IV）
//   - padding：填充方式
//
// 返回值：
//   - []byte：加密后的数据
//   - error：如果加密过程中发生错误，则返回错误信息
func CBCEncrypt(block cipher.Block, src, iv []byte, padding string) ([]byte, error) {
	if block == nil {
		return nil, ErrBlockNil
	}
	blockSize := block.BlockSize()
	src = Padding(padding, src, blockSize)
	encryptData := make([]byte, len(src))
	if len(iv) != blockSize {
		// 自动将IV的长度填充到块大小
		iv = cbcIVPending(iv, blockSize)
		//return nil, errors.New("CBCEncrypt: IV length must equal block size")
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptData, src)

	return encryptData, nil
}

// CBCDecrypt 使用CBC模式解密算法对密文进行解密操作。
// 参数：
//   - block：使用的加密块。
//   - src：待解密的密文。
//   - iv：初始向量，长度必须等于块大小。
//   - padding：填充方式。
//
// 返回值：
//   - []byte：解密后的明文。
//   - error：解密过程中可能发生的错误。
func CBCDecrypt(block cipher.Block, src, iv []byte, padding string) ([]byte, error) {
	if block == nil {
		return nil, ErrBlockNil
	}
	dst := make([]byte, len(src))
	if len(iv) != block.BlockSize() {
		// 自动填充长度到块大小
		iv = cbcIVPending(iv, block.BlockSize())
		//return nil, errors.New("CBCDecrypt: IV length must equal block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, src)
	return UnPadding(padding, dst)
}

// cbcIVPending 根据提供的初始化向量 iv 和区块大小 blockSize 返回一个新的初始化向量。
// 如果 iv 为 nil，则将其设为空字节切片。
// 如果 iv 的长度小于 blockSize，则在 iv 的末尾添加足够的零字节使得长度达到 blockSize。
// 如果 iv 的长度大于 blockSize，则截取 iv 的前 blockSize 个字节作为新的初始化向量。
// 返回新的初始化向量。
func cbcIVPending(iv []byte, blockSize int) []byte {
	if iv == nil {
		iv = []byte{}
	}
	k := len(iv)
	if k < blockSize {
		return append(iv, bytes.Repeat([]byte{0}, blockSize-k)...)
	} else if k > blockSize {
		return iv[0:blockSize]
	}
	return iv
}
