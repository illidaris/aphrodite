package encrypter

import (
	"context"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

// Generate 生成数据密钥DEK
func (c KmsTencentClient) generateDataKey(ctx context.Context, keyId, keySpec string) (string, string, error) {
	request := kms.NewGenerateDataKeyRequest()
	request.KeyId = common.StringPtr(keyId)
	request.KeySpec = common.StringPtr(keySpec)
	rsp, err := c.client.GenerateDataKeyWithContext(ctx, request)
	if err != nil {
		return "", "", err
	}
	return *rsp.Response.Plaintext, *rsp.Response.CiphertextBlob, nil
}
