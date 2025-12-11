package encrypter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"

	"github.com/tjfoc/gmsm/sm4"
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
func TestEncryptStreamSM4(t *testing.T) {
	opts := []Option{
		WithBlockSize(sm4.BlockSize),
		WithCipher(sm4.NewCipher),
		WithEncrypter(cipher.NewCTR),
		WithDecrypter(cipher.NewCTR)}
	secret := []byte("0123456789abcdef")
	raw := "test !@#^&*())_*)_!@#!dataadsdasdasd"
	println("Raw: " + raw)

	enin := bytes.NewBufferString(raw)
	enout := &bytes.Buffer{}
	enerr := EncryptStream(enin, enout, secret, opts...)
	if enerr != nil {
		t.Errorf("EncryptStream failed: %v", enerr)
	}
	println("RawEn: " + base64.StdEncoding.EncodeToString(enout.Bytes()))

	enStr := "if93fMnJBOQENhtyVojPysd8SvuqzmWUDXXSylDL7tOcCwZAKQ+PdBSIQFir4bLJ/rrvRw=="
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

// TestEncryptStream 测试 EncryptStream 函数
func TestEncryptStreamVal(t *testing.T) {
	opts := []Option{
		WithCipher(aes.NewCipher),
		WithEncrypter(cipher.NewCTR),
		WithDecrypter(cipher.NewCTR)}
	secret := []byte("0123456789abcdef")
	raw := "test !@#^&*())_*)_!@#!dataadsdasdasd"
	println("Raw: " + raw)

	enin := bytes.NewBufferString(raw)
	enout := &bytes.Buffer{}
	enerr := EncryptStream(enin, enout, secret, opts...)
	if enerr != nil {
		t.Errorf("EncryptStream failed: %v", enerr)
	}
	println("RawEn: " + base64.StdEncoding.EncodeToString(enout.Bytes()))

	enStr := "uEFtesnTRAuywj7FOrAsjH+UNwkJllePLDdkKyVMh7b3MGdzN3hR+SXI9fV7zHntUaI1Pg=="
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

func TestEncrypt(t *testing.T) {
	opts := []Option{
		WithCipher(sm4.NewCipher),
		WithEncrypter(cipher.NewCTR),
		WithDecrypter(cipher.NewCTR)}
	secret, _ := base64.StdEncoding.DecodeString("5Ye0OiRelFz7FsZx7SWuBw==")

	raw := "test !@#^&*())_*)_!@#!dataadsdasdasd"
	println("Raw: " + raw)

	enin := bytes.NewBufferString(raw)
	enout := &bytes.Buffer{}

	enerr := EncryptStream(enin, enout, secret, opts...)
	if enerr != nil {
		t.Errorf("EncryptStream failed: %v", enerr)
	}

	res := base64.StdEncoding.EncodeToString(enout.Bytes())
	println("RawEn: " + res)
	if res != "uEFtesnTRAuywj7FOrAsjH+UNwkJllePLDdkKyVMh7b3MGdzN3hR+SXI9fV7zHntUaI1Pg==" {
		t.Error("加密失败")
	}
}

func TestDecrypt(t *testing.T) {
	opts := []Option{
		WithCipher(sm4.NewCipher),
		WithEncrypter(cipher.NewCTR),
		WithDecrypter(cipher.NewCTR)}
	secret, _ := base64.StdEncoding.DecodeString("5Ye0OiRelFz7FsZx7SWuBw==")

	enStr := "hKFoLK1RNdE/t2TUrpS+OmFusRlfDdYFIA0dmEWcbYvMaP0U67joXDkzDAQb9HzZ4buzAw=="
	bs, _ := base64.StdEncoding.DecodeString(enStr)
	dein := bytes.NewBuffer(bs)
	deout := &bytes.Buffer{}
	deerr := DecryptStream(dein, deout, secret, opts...)
	if deerr != nil {
		t.Errorf("EncryptStream failed: %v", deerr)
	}
	println("RawDe: " + deout.String())

	if deout.String() != "test !@#^&*())_*)_!@#!dataadsdasdasd" {
		t.Error("解密失败")
	}
}
