/*
https://www.tencentcloud.com/zh/document/product/1030/32774

密钥管理服务 KMS
提供安全合规的密钥的全生命周期管理和数据加解密能力。

对于用户而言，KMS 服务中涉及的核心密钥组件包括:
 1. 用户主密钥 CMK（Customer Master Key，CMK）
    用户主密钥是 KMS 中的核心资源，这些主密钥经过第三方认证硬件安全模块（HSM）的保护，作为用户加密解密的一级密钥。
    KMS 服务主要是针对用户主密钥的管理服务。
    用户主密钥 CMK 是主密钥的逻辑表示。
    CMK 包含元数据，例如密钥 ID、创建日期、描述和密钥状态等。
    通常情况下您可以使用 KMS 的自动生成用户主密钥功能来生成 CMK，同时支持您自有密钥的导入来形成 CMK。
    用户主密钥 CMK 包括用户密钥和云产品密钥两种类型：
    用户密钥是用户通过控制台或 API 来创建的用户主密钥。
    您可以对用户密钥进行创建/启用/禁用/轮换/权限控制等操作。
    云产品密钥是腾讯云产品/服务（例如 CBS、COS、TDSQL 等）在调用密钥管理服务时，自动为用户创建的 CMK。您可以对云产品密钥进行查询及开启密钥轮换操作，不支持禁用、计划删除操作。
 2. 数据加密密钥 DEK（Data Encryption Key，DEK）
    数据加密密钥是基于 CMK 生成的二级密钥，可用于用户本地数据加密解密。
    您可以使用 KMS 用户主密钥（CMK）生成 DEK，但是，KMS 不会存储、管理或跟踪您的 DEK，也不会用于 DEK 执行加密操作。您必须在 KMS 之外使用和管理 DEK。
    一般 DEK 在信封加密流程中使用，通过 DEK 进行本地业务数据的加密。DEK 受用户主密钥 CMK 保护，可以自定义，也可以通过 GenerateDataKey 接口来创建 DEK。

其中 CMK 属于用户的一级密钥，CMK 用于对敏感数据的加解密以及 DEK 的派生。DEK 是信封加密流程中的二级密钥，用于加密业务数据的密钥，受用户主密钥 CMK 的保护。
关于使用 CMK 及 DEK 进行业务加解密的场景，请参见 敏感数据加密 和 信封加密最佳实践。
*/

package encrypter

import (
	"context"
	"encoding/base64"
	"fmt"

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

var _ = IKmsAdapter(KmsTencentClient{})

type KmsTencentClient struct {
	KmsClientOptions
	client *kms.Client
}

func (c KmsTencentClient) GenerateDEK(ctx context.Context, keyId, keySpec string) ([]byte, []byte, error) {
	request := kms.NewGenerateDataKeyRequest()
	request.KeyId = common.StringPtr(keyId)
	request.KeySpec = common.StringPtr(keySpec)
	rsp, err := c.client.GenerateDataKey(request)
	if err != nil {
		return nil, nil, err
	}
	cipherDek := *rsp.Response.CiphertextBlob
	// 解密本地DEK
	plainDek, err := c.DecryptDEK(cipherDek)
	if err != nil {
		return nil, nil, err
	}
	// 解密并校验一次
	if *rsp.Response.Plaintext != base64.StdEncoding.EncodeToString(plainDek) {
		return nil, nil, fmt.Errorf("%s生成DEK不可用", *rsp.Response.RequestId)
	}
	return plainDek, []byte(*rsp.Response.CiphertextBlob), nil
}

// DecryptDEK 使用KMS解密DEK
func (c KmsTencentClient) DecryptDEK(cipherDeK string) ([]byte, error) {
	request := kms.NewDecryptRequest()
	request.CiphertextBlob = &cipherDeK
	rsp, err := c.client.Decrypt(request)
	if err != nil {
		return nil, err
	}
	plain := *rsp.Response.Plaintext
	println("明文：", plain) // TODO
	plainDek, err := base64.StdEncoding.DecodeString(plain)
	if err != nil {
		return nil, err
	}
	return plainDek, nil
}

// DecryptDEK 使用KMS加密DEK
func (c KmsTencentClient) EncryptDEK(keyId string, plaintext []byte) (string, error) {
	encryptRequest := kms.NewEncryptRequest()
	encryptRequest.KeyId = common.StringPtr(keyId)
	encryptRequest.Plaintext = common.StringPtr(base64.StdEncoding.EncodeToString(plaintext))
	encryptResponse, err := c.client.Encrypt(encryptRequest)
	if err != nil {
		return "", err
	}
	return *encryptResponse.Response.CiphertextBlob, nil
}
