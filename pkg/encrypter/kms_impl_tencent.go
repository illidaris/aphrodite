package encrypter

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

func NewKmsTencent(opts ...KmsClientOption) (*KmsTencentClient, error) {
	c := &KmsTencentClient{}
	for _, opt := range opts {
		opt(&c.KmsClientOptions)
	}
	credential := common.NewCredential(c.AppId, c.Secret)
	cpf := profile.NewClientProfile()
	client, err := kms.NewClient(credential, c.Region, cpf)
	if err != nil {
		return c, err
	}
	c.client = client
	return c, nil
}

type KmsTencentClient struct {
	KmsClientOptions
	client *kms.Client
	store  IKmsStore
}

func (c KmsTencentClient) GenerateDEK(ctx context.Context, opts ...KmsOption) error {
	option := newKmsOptions(opts...)

	oldCipherDekBs, err := c.store.Get(ctx, option.KeyId)
	if err != nil {
		return err
	}
	if len(oldCipherDekBs) > 0 {
		return errors.New("DEK已经存在")
	}

	request := kms.NewGenerateDataKeyRequest()
	request.KeyId = common.StringPtr(option.KeyId)
	request.KeySpec = common.StringPtr(option.KeySpec)

	rsp, err := c.client.GenerateDataKey(request)
	if err != nil {
		return err
	}

	cipherDekBs, err := base64.StdEncoding.DecodeString(*rsp.Response.CiphertextBlob)
	if err != nil {
		return err
	}
	// 解密本地DEK
	plainDek, err := c.DecryptDEK(cipherDekBs)
	if err != nil {
		return err
	}
	// 解密并校验一次
	if *rsp.Response.Plaintext != base64.StdEncoding.EncodeToString(plainDek) {
		return fmt.Errorf("%s生成DEK不可用", *rsp.Response.RequestId)
	}
	affect, err := c.store.Save(ctx, option.KeyId, cipherDekBs)
	if err != nil {
		return err
	}
	if affect == 0 {
		return errors.New("保存DEK失败")
	}
	return nil
}
func (c KmsTencentClient) Encrypt(val []byte, opts ...KmsOption) ([]byte, error) {
	return c.EncryptCtx(context.Background(), val, opts...)
}

func (c KmsTencentClient) Decrypt(val []byte, opts ...KmsOption) ([]byte, error) {
	return c.DecryptCtx(context.Background(), val, opts...)
}

func (c KmsTencentClient) EncryptCtx(ctx context.Context, val []byte, opts ...KmsOption) ([]byte, error) {
	outBs := []byte{}
	in := bytes.NewBuffer(val)
	out := bytes.NewBuffer(outBs)
	err := c.EncryptStream(ctx, in, out, opts...)
	if err != nil {
		return outBs, err
	}
	return outBs, nil
}

func (c KmsTencentClient) EncryptStream(ctx context.Context, in io.Reader, out io.Writer, opts ...KmsOption) error {
	option := newKmsOptions(opts...)
	// 提取数据库中保存的cipherDek
	cipherDekBs, err := c.store.Get(ctx, option.KeyId)
	if err != nil {
		return err
	}
	// 解密本地DEK
	plainDek, err := c.DecryptDEK(cipherDekBs)
	if err != nil {
		return err
	}
	// plainDek DEK明文用户缓存在内存中使用，对数据进行本地加密
	err = EncryptStream(in, out, plainDek, option.AESOption...)
	if err != nil {
		return err
	}
	return nil
}

func (c KmsTencentClient) DecryptCtx(ctx context.Context, val []byte, opts ...KmsOption) ([]byte, error) {
	outBs := []byte{}
	in := bytes.NewBuffer(val)
	out := bytes.NewBuffer(outBs)
	err := c.DecryptStream(ctx, in, out, opts...)
	if err != nil {
		return outBs, err
	}
	return outBs, nil
}

func (c KmsTencentClient) DecryptStream(ctx context.Context, in io.Reader, out io.Writer, opts ...KmsOption) error {
	option := newKmsOptions(opts...)
	// 提取数据库中保存的cipherDek
	cipherDekBs, err := c.store.Get(ctx, option.KeyId)
	if err != nil {
		return err
	}
	// 解密本地DEK
	plainDek, err := c.DecryptDEK(cipherDekBs)
	if err != nil {
		return err
	}
	// plainDek DEK明文用户缓存在内存中使用，对数据进行本地加密
	err = DecryptStream(in, out, plainDek, option.AESOption...)
	if err != nil {
		return err
	}
	return nil
}

// DecryptDEK 使用KMS解密DEK
func (c KmsTencentClient) DecryptDEK(cipherDeKBs []byte) ([]byte, error) {
	request := kms.NewDecryptRequest()
	cipherDeK := base64.StdEncoding.EncodeToString(cipherDeKBs)
	request.CiphertextBlob = &cipherDeK
	rsp, err := c.client.Decrypt(request)
	if err != nil {
		return nil, err
	}
	plainDek, err := base64.StdEncoding.DecodeString(*rsp.Response.Plaintext)
	if err != nil {
		return nil, err
	}
	return plainDek, nil
}

// DecryptDEK 使用KMS加密DEK
func (c KmsTencentClient) EncryptDEK(keyId string, plaintext []byte) ([]byte, error) {
	encryptRequest := kms.NewEncryptRequest()
	encryptRequest.KeyId = common.StringPtr(keyId)
	encryptRequest.Plaintext = common.StringPtr(base64.StdEncoding.EncodeToString(plaintext))
	encryptResponse, err := c.client.Encrypt(encryptRequest)
	if err != nil {
		return nil, err
	}
	raw := *encryptResponse.Response.CiphertextBlob
	return base64.StdEncoding.DecodeString(raw)
}
