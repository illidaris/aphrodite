package encrypter

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
)

type IKms interface {
	Encrypt(val []byte, opts ...KmsOption) ([]byte, error)
	EncryptCtx(ctx context.Context, val []byte, opts ...KmsOption) ([]byte, error)
	Decrypt(val []byte, opts ...KmsOption) ([]byte, error)
	DecryptCtx(ctx context.Context, val []byte, opts ...KmsOption) ([]byte, error)
}

type IKmsStore interface {
	Save(ctx context.Context, key string, val []byte) (int64, error)
	Get(ctx context.Context, key string) ([]byte, error)
}

type KmsClientOption func(*KmsClientOptions)

type KmsClientOptions struct {
	Region string
	AppId  string
	Secret string
}

type KmsOption func(*KmsOptions)

func newKmsOptions(opts ...KmsOption) *KmsOptions {
	options := &KmsOptions{
		KeySpec: "AES_128",
		AESOption: []Option{
			WithCipher(aes.NewCipher),
			WithDecrypter(cipher.NewCTR),
			WithDecrypter(cipher.NewCTR),
		},
	}
	for _, v := range opts {
		v(options)
	}
	return options
}

type KmsOptions struct {
	KeyId     string
	KeySpec   string
	AESOption []Option
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
