package oauth2

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"

	"github.com/illidaris/aphrodite/pkg/encrypter"
	"github.com/tjfoc/gmsm/sm4"
)

var (
	aesOpts = []encrypter.Option{
		encrypter.WithBlockSize(sm4.BlockSize),
		encrypter.WithCipher(sm4.NewCipher),
		encrypter.WithEncrypter(cipher.NewCTR),
		encrypter.WithDecrypter(cipher.NewCTR),
	}
)

func AESEncode[T any](v T, secret string) (string, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	deOutBs := bytes.NewBuffer(nil)
	inBs := bytes.NewBuffer(bs)
	err = encrypter.EncryptStream(
		inBs, deOutBs,
		[]byte(secret),
		aesOpts...,
	)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(deOutBs.Bytes()), nil
}

func AESDecode[T any](res *T, str, secret string) error {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	deOutBs := bytes.NewBuffer(nil)
	inBs := bytes.NewBuffer(decodeBytes)
	err = encrypter.DecryptStream(
		inBs, deOutBs,
		[]byte(secret),
		aesOpts...,
	)
	if err != nil {
		return err
	}
	err = json.Unmarshal(deOutBs.Bytes(), res)
	if err != nil {
		return err
	}
	return nil
}
