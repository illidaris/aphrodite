package encrypter

import (
	"context"
	"crypto/cipher"
	"fmt"
	"io"
	"sync"

	"github.com/tjfoc/gmsm/sm4"
)

type KmsClientOption func(*KmsClientOptions)

type KmsClientOptions struct {
	Region string
	AppId  string
	Secret string
}

func WithKmsClientAppId(v string) KmsClientOption {
	return func(o *KmsClientOptions) {
		o.AppId = v
	}
}

func WithKmsClientSecret(v string) KmsClientOption {
	return func(o *KmsClientOptions) {
		o.Secret = v
	}
}

func WithKmsClientRegion(v string) KmsClientOption {
	return func(o *KmsClientOptions) {
		o.Region = v
	}
}

type KmsOption func(*KmsOptions)

func newKmsOptions(opts ...KmsOption) *KmsOptions {
	options := &KmsOptions{
		Id:      SECRET_DEF,
		KeySpec: SPEC_KEY_AES_128,
		AESOption: []Option{
			WithCipher(sm4.NewCipher),
			WithBlockSize(sm4.BlockSize),
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
	Id                string
	KeyId             string
	KeySpec           string
	AESOption         []Option
	DecryptStreamFunc Encryptfunc
	EncryptStreamFunc Encryptfunc
}

func WithKmsId(v string) KmsOption {
	return func(o *KmsOptions) {
		o.Id = v
	}
}
func WithKmsKeyId(keyId string) KmsOption {
	return func(o *KmsOptions) {
		o.KeyId = keyId
	}
}
func WithKmsAESOption(vs ...Option) KmsOption {
	return func(o *KmsOptions) {
		o.AESOption = append(o.AESOption, vs...)
	}
}

func WithKmsEncryptStreamFunc(v Encryptfunc) KmsOption {
	return func(o *KmsOptions) {
		o.EncryptStreamFunc = v
	}
}

func WithKmsDecryptStreamFunc(v Encryptfunc) KmsOption {
	return func(o *KmsOptions) {
		o.DecryptStreamFunc = v
	}
}

func newEmbeddedCache() *embeddedCache {
	return &embeddedCache{
		m: &sync.Map{},
	}
}

var _ = IKmsCache(embeddedCache{})

type embeddedCache struct {
	m *sync.Map
}

func (i embeddedCache) DekPlainSave(ctx context.Context, dek *DEKPlainEntry) (int64, error) {
	i.m.Store(dek.Id, dek)
	println("明文：", dek.Plain) // TODO
	return 1, nil
}
func (i embeddedCache) DekPlainGet(ctx context.Context, id string) (*DEKPlainEntry, error) {
	v, ok := i.m.Load(id)
	if !ok {
		return nil, fmt.Errorf("no found %v", id)
	}
	return v.(*DEKPlainEntry), nil
}
