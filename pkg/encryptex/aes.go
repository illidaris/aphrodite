package openssl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AesECBEncrypt 使用AES算法的ECB模式进行加密。
// 参数src为待加密的数据，key为加密密钥，padding为填充方式。
// 返回加密后的数据和可能的错误。
func AesECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := AesNewCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, padding)
}

// AesECBDecrypt 使用AES算法的ECB模式进行解密
// 参数：
//   - src: 待解密的数据
//   - key: 解密密钥
//   - padding: 填充方式
//
// 返回值：
//   - []byte: 解密后的数据
//   - error: 解密过程中遇到的错误
func AesECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := AesNewCipher(key)
	if err != nil {
		return nil, err
	}

	return ECBDecrypt(block, src, padding)
}

// AesCBCEncrypt 使用AES算法进行CBC模式加密。
// 参数src是要加密的数据。
// 参数key是加密密钥。
// 参数iv是初始化向量。
// 参数padding是填充方式。
// 返回加密后的数据和可能的错误。
func AesCBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := AesNewCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCEncrypt(block, src, iv, padding)
}

// AesCBCDecrypt 使用AES算法的CBC模式进行密码解密。
// 参数src为待解密的数据，key为解密密钥，iv为初始向量，padding为填充方式。
// 返回解密后的数据和可能的错误。
func AesCBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := AesNewCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCDecrypt(block, src, iv, padding)
}

// AesNewCipher 使用AES算法创建一个新的加密器。
// 参数key为加密使用的密钥。
// 返回值为一个实现cipher.Block接口的加密器和一个错误。
func AesNewCipher(key []byte) (cipher.Block, error) {
	return aes.NewCipher(aesKeyPending(key))
}

// aesKeyPending 根据给定的密钥生成待用的密钥。
// 如果密钥长度小于等于16，则返回一个长度为16的密钥。
// 如果密钥长度大于16且小于等于24，则返回一个长度为24的密钥。
// 如果密钥长度大于24且小于等于32，则返回一个长度为32的密钥。
// 如果密钥长度大于32，则返回一个长度为32的密钥。
// 如果密钥长度大于16且小于等于32，则在密钥末尾添加0填充至32。
func aesKeyPending(key []byte) []byte {
	k := len(key)
	var count int
	switch true {
	case k <= 16:
		count = 16 - k
	case k <= 24:
		count = 24 - k
	case k <= 32:
		count = 32 - k
	default:
		return key[:32]
	}
	if count == 0 {
		return key
	}
	return append(key, bytes.Repeat([]byte{0}, count)...)
}
