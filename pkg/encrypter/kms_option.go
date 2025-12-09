package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type KmsClientOption func(*KmsClientOptions)

type KmsClientOptions struct {
	Region string
	AppId  string
	Secret string
}

type KmsOption func(*KmsOptions)

func newKmsOptions(opts ...KmsOption) *KmsOptions {
	options := &KmsOptions{
		KeySpec: SPEC_KEY_AES_128,
		AESOption: []Option{
			WithCipher(aes.NewCipher),
			WithDecrypter(cipher.NewCTR),
			WithDecrypter(cipher.NewCTR),
		},
		EncryptStreamFunc: EncryptStream,
		DecryptStreamFunc: DecryptStream,
	}
	for _, v := range opts {
		v(options)
	}
	return options
}

type Encryptfunc func(in io.Reader, out io.Writer, secret []byte, opts ...Option) error
type KmsOptions struct {
	KeyId             string
	KeySpec           string
	AESOption         []Option
	DecryptStreamFunc Encryptfunc
	EncryptStreamFunc Encryptfunc
}

func WithKeyId(keyId string) KmsOption {
	return func(o *KmsOptions) {
		o.KeyId = keyId
	}
}
func WithAESOption(vs ...Option) KmsOption {
	return func(o *KmsOptions) {
		o.AESOption = append(o.AESOption, vs...)
	}
}

func WithEncryptStreamFunc(v Encryptfunc) KmsOption {
	return func(o *KmsOptions) {
		o.EncryptStreamFunc = v
	}
}

func WithDecryptStreamFunc(v Encryptfunc) KmsOption {
	return func(o *KmsOptions) {
		o.DecryptStreamFunc = v
	}
}
