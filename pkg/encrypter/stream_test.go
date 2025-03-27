package encrypter

import (
	"bytes"
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
