package encrypter

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type dataEncryptOption struct {
	Secret []byte
}

type DataEncryptOptionFunc func(*dataEncryptOption)

func WithSecret(secret []byte) DataEncryptOptionFunc {
	return func(o *dataEncryptOption) {
		o.Secret = secret
	}
}

func WithSecretString(secret string) DataEncryptOptionFunc {
	return func(o *dataEncryptOption) {
		o.Secret = []byte(secret)
	}
}
func newDataEncryptOption(opts ...DataEncryptOptionFunc) *dataEncryptOption {
	opt := &dataEncryptOption{}
	for _, f := range opts {
		f(opt)
	}
	return opt
}
func DataEncrypt(value []byte, opts ...DataEncryptOptionFunc) ([]byte, error) {
	opt := newDataEncryptOption(opts...)
	if len(opt.Secret) == 0 {
		return value, errors.New("secret is nil")
	}
	outBs := &bytes.Buffer{}
	err := EncryptStream(bytes.NewReader(value), outBs, opt.Secret,
		WithDecrypter(cipher.NewCTR),
		WithEncrypter(cipher.NewCTR),
	)
	if err != nil {
		return outBs.Bytes(), err
	}
	return outBs.Bytes(), nil
}

func DataDecrypt(value []byte, opts ...DataEncryptOptionFunc) ([]byte, error) {
	opt := newDataEncryptOption(opts...)
	if len(opt.Secret) == 0 {
		return value, errors.New("secret is nil")
	}
	outBs := &bytes.Buffer{}
	err := DecryptStream(bytes.NewReader(value), outBs, opt.Secret,
		WithDecrypter(cipher.NewCTR),
		WithEncrypter(cipher.NewCTR),
	)
	if err != nil {
		return outBs.Bytes(), err
	}
	return outBs.Bytes(), nil
}

func StringEncrypt(value string, opts ...DataEncryptOptionFunc) (string, error) {
	bs, err := DataEncrypt([]byte(value), opts...)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bs), nil
}

func StringDecrypt(value string, opts ...DataEncryptOptionFunc) (string, error) {
	rawBs, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	bs, err := DataDecrypt(rawBs, opts...)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func JsonEncrypt[T any](value T, opts ...DataEncryptOptionFunc) (string, error) {
	rawBs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	bs, err := DataEncrypt(rawBs, opts...)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bs), nil
}

func JsonDecrypt[T any](value string, opts ...DataEncryptOptionFunc) (T, error) {
	var res T
	rawBs, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return res, err
	}
	bs, err := DataDecrypt(rawBs, opts...)
	if err != nil {
		return res, err
	}
	_ = json.Unmarshal(bs, &res)
	return res, nil
}
