package openssl

import (
	"crypto/cipher"
)

// ECBEncrypt 使用ECB模式加密算法对给定的数据进行加密。
// 参数：
//   - block：加密算法的块
//   - src：要加密的数据
//   - padding：填充方式
//
// 返回值：
//   - []byte：加密后的数据
//   - error：如果加密算法的块为nil，则返回ErrBlockNil错误
func ECBEncrypt(block cipher.Block, src []byte, padding string) ([]byte, error) {
	if block == nil {
		return nil, ErrBlockNil
	}
	blockSize := block.BlockSize()
	src = Padding(padding, src, blockSize)

	encryptData := make([]byte, len(src))

	ecb := NewECBEncrypter(block)
	ecb.CryptBlocks(encryptData, src)

	return encryptData, nil
}

// ECBDecrypt 使用ECB模式解密算法对给定的数据进行解密操作。
// 参数：
//   - block：加密算法的块密码。
//   - src：待解密的数据。
//   - padding：填充方式。
//
// 返回值：
//   - []byte：解密后的数据。
//   - error：解密操作过程中遇到的错误。
func ECBDecrypt(block cipher.Block, src []byte, padding string) ([]byte, error) {
	dst := make([]byte, len(src))

	// 创建ECB解密器
	mode := NewECBDecrypter(block)
	// 对数据进行解密操作
	mode.CryptBlocks(dst, src)

	// 对解密后的数据进行填充操作
	return UnPadding(padding, dst)
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

// newECB 返回一个初始化的 ecb 结构体指针。
// 参数 b 为密码块，ecb 结构体的 b 字段与参数 b 相同。
// 如果 b 不为 nil，则将密码块的块大小赋值给 ecb 结构体的 blockSize 字段。
// 返回初始化的 ecb 结构体指针。
func newECB(b cipher.Block) *ecb {
	res := &ecb{
		b: b,
	}
	if b != nil {
		res.blockSize = b.BlockSize()
	}
	return res
}

type ecbEncrypter ecb

// NewECBEncrypter 返回一个使用ECB模式的加密器。
// 参数b为一个cipher.Block类型的参数，用于初始化加密器。
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter 返回一个ECB解密器的实例。
// 参数b为一个cipher.Block，用于初始化解密器。
// 返回一个cipher.BlockMode类型的指针。
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
