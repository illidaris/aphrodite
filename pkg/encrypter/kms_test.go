package encrypter

import (
	"bytes"
	"context"
	"encoding/base64"
	"testing"
)

var _ = IKmsStore(&printStore{})

type printStore struct {
	m map[string]*DEKEntry
}

func (s printStore) DekSave(ctx context.Context, dek *DEKEntry) (int64, error) {
	println(dek.Id)
	s.m[dek.Id] = dek
	println(dek.Cipher)
	return 1, nil
}
func (s printStore) DekFind(ctx context.Context, ids ...string) ([]*DEKEntry, error) {
	ds := []*DEKEntry{}
	for _, id := range ids {
		ds = append(ds, s.m[id])
	}
	return ds, nil
}
func TestKms(t *testing.T) {
	s := &printStore{
		m: map[string]*DEKEntry{},
	}
	c, err := NewKmEmbed(
		WithKmsClientAppId("AK*********"),
		WithKmsClientSecret("asdasdasd"),
		WithKmsClientRegion("*******"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	manage := NewKmsManage(c, s, nil)

	ctx := context.Background()
	keyIId := "asdasdas"
	_, err = manage.Generate(ctx, WithKmsKeyId(keyIId))
	if err != nil {
		t.Error(err)
		return
	}

	raw := "test !@#^&*())_*)_!@#!dataadsdassada 谁呀dasd"
	println("Raw: " + raw)

	enin := bytes.NewBufferString(raw)
	enout := &bytes.Buffer{}
	enerr := manage.Encrypt(ctx, enin, enout)
	if enerr != nil {
		t.Errorf("EncryptStream failed: %v", enerr)
	}
	enStr := base64.StdEncoding.EncodeToString(enout.Bytes())
	println("RawEn: " + enStr)

	bs, _ := base64.StdEncoding.DecodeString(enStr)
	dein := bytes.NewBuffer(bs)
	deout := &bytes.Buffer{}
	deerr := manage.Decrypt(ctx, dein, deout)
	if deerr != nil {
		t.Errorf("EncryptStream failed: %v", deerr)
	}
	println("RawDe: " + deout.String())

	if deout.String() != raw {
		t.Error("解密失败")
	}
}
