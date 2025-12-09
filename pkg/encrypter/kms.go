package encrypter

import (
	"context"
)

const (
	SPEC_KEY_AES_128 = "AES_128"
	SPEC_KEY_AES_256 = "AES_256"
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
