package encrypter

import (
	"context"

	"github.com/illidaris/rest/signature"
)

func NewKmEmbed(opts ...KmsClientOption) (*KmEmbedClient, error) {
	c := &KmEmbedClient{}
	for _, opt := range opts {
		opt(&c.KmsClientOptions)
	}
	return c, nil
}

var _ = IKmsAdapter(KmEmbedClient{})

type KmEmbedClient struct {
	KmsClientOptions
}

func (c KmEmbedClient) GenerateDEK(ctx context.Context, keyId, keySpec string) ([]byte, []byte, error) {
	raw := signature.RandString(16, []rune("abcdefghijklmnopqrstuvwxyz_-!@#$%^&*()"))
	return []byte(raw), []byte(raw), nil
}

// DecryptDEK 使用KMS解密DEK
func (c KmEmbedClient) DecryptDEK(cipherDeK string) ([]byte, error) {
	return []byte(cipherDeK), nil
}

// DecryptDEK 使用KMS加密DEK
func (c KmEmbedClient) EncryptDEK(_ string, plaintext []byte) (string, error) {
	return string(plaintext), nil
}
