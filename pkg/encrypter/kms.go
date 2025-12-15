package encrypter

import (
	"context"
)

const (
	SECRET_DEF       = "_apsecret_default"
	SPEC_KEY_AES_128 = "AES_128"
	SPEC_KEY_AES_256 = "AES_256"
)

type IKmsAdapter interface {
	GenerateDEK(ctx context.Context, keyId, keySpec string) ([]byte, []byte, error)
	DecryptDEK(cipherDeK string) ([]byte, error)
	EncryptDEK(keyId string, plaintext []byte) (string, error)
}

type IKmsStore interface {
	DekSave(ctx context.Context, dek *DEKEntry) (int64, error)
	DekFind(ctx context.Context, ids ...string) ([]*DEKEntry, error)
}

type IKmsCache interface {
	DekPlainSave(ctx context.Context, dek *DEKPlainEntry) (int64, error)
	DekPlainGet(ctx context.Context, id string) (*DEKPlainEntry, error)
}

type DEKEntry struct {
	Id       string `json:"id"`       // id
	Name     string `json:"name"`     // 名称
	KmsType  int32  `json:"kmsType"`  // kms类型
	KeyId    string `json:"keyId"`    // key Id
	Cipher   string `json:"cipher"`   // 加密DEK
	CreateAt int64  `json:"createAt"` // 生成时间
	Describe string `json:"describe"` // 描述
}

func (i DEKEntry) WithPlain(v string) *DEKPlainEntry {
	return &DEKPlainEntry{
		DEKEntry: i,
		Plain:    v,
	}
}

type DEKPlainEntry struct {
	DEKEntry
	Plain string `json:"plain"` // 明文 只能在内存中，不能落地
}
