package openssl

import "crypto/des"

// Des3ECBEncrypt 使用Triple DES算法进行ECB模式加密。
// 参数src是要加密的数据。
// 参数key是加密密钥。
// 参数padding是填充方式，目前仅支持PKCS5Padding。
// 返回加密后的数据和可能的错误。
func Des3ECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, padding)
}

// Des3ECBDecrypt 使用TripleDESCipher算法进行密码解密
func Des3ECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	// 创建TripleDESCipher密码块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	// 调用ECBDecrypt函数进行密码解密
	return ECBDecrypt(block, src, padding)
}

// Des3CBCEncrypt 使用Triple DES算法和CBC模式对给定的源数据进行加密。
// 参数src是要加密的数据。
// 参数key是加密密钥。
// 参数iv是初始化向量。
// 参数padding是填充模式。
// 返回加密后的数据和可能的错误。
func Des3CBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCEncrypt(block, src, iv, padding)
}

// Des3CBCDecrypt 使用Triple DES算法和CBC模式进行密码解密。
// 参数src为待解密的数据，key为解密密钥，iv为初始向量，padding为填充方式。
// 返回解密后的数据和错误信息。
func Des3CBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	return CBCDecrypt(block, src, iv, padding)
}
