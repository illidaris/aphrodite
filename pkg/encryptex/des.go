package openssl

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// DesECBEncrypt 使用DES算法的ECB模式进行加密。
// 参数src是要加密的数据。
// 参数key是加密密钥。
// 参数padding是填充方式，可选值为"PKCS5Padding"或"PKCS7Padding"。
// 返回加密后的数据和可能的错误。
func DesECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := DesNewCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, padding)
}

// DesECBDecrypt 使用DES算法的ECB模式进行解密。
// 参数src为待解密的数据，key为解密密钥，padding为填充方式。
// 返回解密后的数据和错误信息（如果有）。
func DesECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	// 创建DES加密算法的实例
	block, err := DesNewCipher(key)
	if err != nil {
		return nil, err
	}

	// 调用ECBDecrypt函数进行解密
	return ECBDecrypt(block, src, padding)
}

// DesCBCEncrypt 使用DES算法进行CBC模式加密。
// 参数：
//   - src：待加密的数据
//   - key：加密密钥
//   - iv：初始化向量
//   - padding：填充方式
//
// 返回值：
//   - []byte：加密后的数据
//   - error：如果加密过程中发生错误，返回错误信息
func DesCBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := DesNewCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCEncrypt(block, src, iv, padding)
}

// DesCBCDecrypt 使用DES算法的CBC模式进行密码解密
// 参数：
//   - src: 待解密的数据
//   - key: 解密密钥
//   - iv: 初始化向量
//   - padding: 填充方式
//
// 返回值：
//   - []byte: 解密后的数据
//   - error: 如果解密失败，返回错误信息
func DesCBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := DesNewCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCDecrypt(block, src, iv, padding)
}

// DesNewCipher 用于创建一个DES加密算法的密码块。
// 参数 key 是加密算法的密钥，长度必须为8字节。
// 如果 key 的长度小于8字节，则在 key 的末尾补0直到长度为8字节。
// 如果 key 的长度大于8字节，则截取前8字节作为密钥。
// 返回一个 cipher.Block 对象和一个 error 对象。
// 如果密钥长度不正确，将返回一个错误。
func DesNewCipher(key []byte) (cipher.Block, error) {
	if len(key) < 8 {
		key = append(key, bytes.Repeat([]byte{0}, 8-len(key))...)
	} else if len(key) > 8 {
		key = key[:8]
	}

	return des.NewCipher(key)
}
