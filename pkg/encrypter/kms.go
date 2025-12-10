package encrypter

import (
	"context"
)

const (
	SPEC_KEY_AES_128 = "AES_128"
	SPEC_KEY_AES_256 = "AES_256"
)

type IKms interface {
	GenerateDEK(ctx context.Context, opts ...KmsOption) error
	Encrypt(val []byte, opts ...KmsOption) ([]byte, error)
	EncryptCtx(ctx context.Context, val []byte, opts ...KmsOption) ([]byte, error)
	Decrypt(val []byte, opts ...KmsOption) ([]byte, error)
	DecryptCtx(ctx context.Context, val []byte, opts ...KmsOption) ([]byte, error)
}

type IKmsStore interface {
	DekSave(ctx context.Context, key string, val string) (int64, error)
	DekGet(ctx context.Context, key string) (string, error)
}
