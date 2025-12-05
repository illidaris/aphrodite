package encrypter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"
)

// TestEncryptStream 测试 EncryptStream 函数
func TestEncryptStream(t *testing.T) {
	secret := []byte("0123456789abcdef")
	raw := "test data"

	println("Raw: " + raw)

	enin := bytes.NewBufferString(raw)
	enout := &bytes.Buffer{}
	enerr := EncryptStream(enin, enout, secret)
	if enerr != nil {
		t.Errorf("EncryptStream failed: %v", enerr)
	}
	println("RawEn: " + enout.String())

	dein := bytes.NewBufferString(enout.String())
	deout := &bytes.Buffer{}
	deerr := DecryptStream(dein, deout, secret)
	if deerr != nil {
		t.Errorf("EncryptStream failed: %v", deerr)
	}
	println("RawDe: " + deout.String())

	if deout.String() != raw {
		t.Error("解密失败")
	}
}

// TestEncryptStream 测试 EncryptStream 函数
func TestEncryptStreamVal(t *testing.T) {
	opts := []Option{
		WithCipher(aes.NewCipher),
		WithEncrypter(cipher.NewCTR),
		WithDecrypter(cipher.NewCTR)}
	secret := []byte("0123456789abcdef")
	raw := "test data"
	println("Raw: " + raw)

	enin := bytes.NewBufferString(raw)
	enout := &bytes.Buffer{}
	enerr := EncryptStream(enin, enout, secret, opts...)
	if enerr != nil {
		t.Errorf("EncryptStream failed: %v", enerr)
	}
	println("RawEn: " + base64.StdEncoding.EncodeToString(enout.Bytes()))

	enStr := "EkG5Po44Q4ti9H48LK7Qkw=="
	bs, _ := base64.StdEncoding.DecodeString(enStr)
	dein := bytes.NewBuffer(bs)
	deout := &bytes.Buffer{}
	deerr := DecryptStream(dein, deout, secret, opts...)
	if deerr != nil {
		t.Errorf("EncryptStream failed: %v", deerr)
	}
	println("RawDe: " + deout.String())

	if deout.String() != raw {
		t.Error("解密失败")
	}
}
