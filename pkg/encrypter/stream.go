package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// newOptions 创建一个新的 options 对象，并应用传入的 Option 函数进行配置。
// opts: 可选的配置函数，用于修改 options 对象的默认值。
// 返回值: 配置后的 options 对象指针。
func newOptions(opts ...Option) *options {
	o := &options{
		Cipher:    aes.NewCipher,
		Encrypter: cipher.NewCFBEncrypter,
		Decrypter: cipher.NewCFBDecrypter,
	}
	// 应用所有传入的 Option 函数，修改 options 对象的配置
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// options 结构体用于存储加密和解密的配置选项。
type options struct {
	Cipher    func([]byte) (cipher.Block, error)       // 用于生成加密块的函数
	Encrypter func(cipher.Block, []byte) cipher.Stream // 用于创建加密流的函数
	Decrypter func(cipher.Block, []byte) cipher.Stream // 用于创建解密流的函数
}

// Option 是一个函数类型，用于修改 options 对象的配置。
type Option func(*options)

// WithCipher 返回一个 Option 函数，用于设置 options 中的 Cipher 字段。
// cipher: 用于生成加密块的函数。
// 返回值: 一个 Option 函数，用于修改 options 对象的 Cipher 字段。
func WithCipher(cipher func([]byte) (cipher.Block, error)) Option {
	return func(o *options) {
		o.Cipher = cipher
	}
}

// WithEncrypter 返回一个 Option 函数，用于设置 options 中的 Encrypter 字段。
// encrypter: 用于创建加密流的函数。
// 返回值: 一个 Option 函数，用于修改 options 对象的 Encrypter 字段。
func WithEncrypter(encrypter func(cipher.Block, []byte) cipher.Stream) Option {
	return func(o *options) {
		o.Encrypter = encrypter
	}
}

// WithDecrypter 返回一个 Option 函数，用于设置 options 中的 Decrypter 字段。
// decrypter: 用于创建解密流的函数。
// 返回值: 一个 Option 函数，用于修改 options 对象的 Decrypter 字段。
func WithDecrypter(decrypter func(cipher.Block, []byte) cipher.Stream) Option {
	return func(o *options) {
		o.Decrypter = decrypter
	}
}

// EncryptStream 对输入流进行加密，并将结果写入输出流。
// in: 输入流，包含待加密的数据。
// out: 输出流，用于写入加密后的数据。
// secret: 加密密钥。
// opts: 可选的配置函数，用于修改加密选项。
// 返回值: 如果加密过程中发生错误，则返回错误信息；否则返回 nil。
func EncryptStream(in io.Reader, out io.Writer, secret []byte, opts ...Option) error {
	// 加载配置
	o := newOptions(opts...)
	// 生成密钥和初始化向量
	block, err := o.Cipher(secret)
	if err != nil {
		return err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	// 写入初始化向量
	if _, err := out.Write(iv); err != nil {
		return err
	}
	// 创建加密流
	stream := o.Encrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: out}
	// 将输入流的数据加密并写入输出流
	if _, err := io.Copy(writer, in); err != nil {
		return err
	}
	return nil
}

// DecryptStream 对输入流进行解密，并将结果写入输出流。
// in: 输入流，包含待解密的数据。
// out: 输出流，用于写入解密后的数据。
// secret: 解密密钥。
// opts: 可选的配置函数，用于修改解密选项。
// 返回值: 如果解密过程中发生错误，则返回错误信息；否则返回 nil。
func DecryptStream(in io.Reader, out io.Writer, secret []byte, opts ...Option) error {
	// 加载配置
	o := newOptions(opts...)
	// 生成密钥和初始化向量
	block, err := o.Cipher(secret)
	if err != nil {
		return err
	}
	// 读取初始化向量
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(in, iv); err != nil {
		return err
	}
	// 创建解密流
	stream := o.Decrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: in}
	// 将输入流的数据解密并写入输出流
	if _, err := io.Copy(out, reader); err != nil {
		return err
	}
	return nil
}
