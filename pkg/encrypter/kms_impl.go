package encrypter

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"time"
)

func NewKmsManage(kms IKmsAdapter, store IKmsStore, cache IKmsCache) *KmsManage {
	if cache == nil {
		cache = newEmbeddedCache()
	}
	m := &KmsManage{
		core:  kms,
		store: store,
		cache: cache,
	}
	return m
}

type KmsManage struct {
	core  IKmsAdapter
	store IKmsStore
	cache IKmsCache
}

// 生成DEK， kms server 生成 落地 数据库
func (i KmsManage) Generate(ctx context.Context, opts ...KmsOption) (*DEKPlainEntry, error) {
	option := newKmsOptions(opts...)
	dek := &DEKEntry{}
	dek.Id = option.Id
	dek.KeyId = option.KeyId
	dek.CreateAt = time.Now().Unix()
	plainDekBs, cipherDekBs, err := i.core.GenerateDEK(ctx, option.KeyId, option.KeySpec)
	if err != nil {
		return nil, err
	}
	dek.Cipher = base64.RawStdEncoding.EncodeToString(cipherDekBs)
	// 落库
	storeAffect, storeErr := i.store.DekSave(ctx, dek)
	if storeErr != nil {
		return nil, storeErr
	}
	if storeAffect == 0 {
		return nil, errors.New("generate store affect 0")
	}
	plain := base64.StdEncoding.EncodeToString(plainDekBs)
	dekPlain := dek.WithPlain(plain)
	// 落户缓存
	_, cacheErr := i.cache.DekPlainSave(ctx, dekPlain)
	if cacheErr != nil {
		return dekPlain, cacheErr
	}
	return dekPlain, nil
}

func (i KmsManage) LoadByIds(ctx context.Context, ids ...string) ([]*DEKPlainEntry, error) {
	res := []*DEKPlainEntry{}
	deks, err := i.store.DekFind(ctx, ids...)
	if err != nil {
		return res, err
	}
	for _, v := range deks {
		dekPlain, err := i.DekPlain(ctx, v)
		if err != nil {
			continue
		}
		affect, err := i.cache.DekPlainSave(ctx, dekPlain)
		if err != nil {
			continue
		}
		if affect == 0 {
			continue
		}
		res = append(res, dekPlain)
	}
	return res, nil
}

func (i KmsManage) DekPlain(ctx context.Context, cipherDek *DEKEntry) (*DEKPlainEntry, error) {
	// 通过KMS远程服务解密
	plainnDekBs, err := i.core.DecryptDEK(cipherDek.Cipher)
	if err != nil {
		return nil, err
	}
	// 生成具备明文的密钥数据
	plainDekBs64 := base64.StdEncoding.EncodeToString(plainnDekBs)
	res := cipherDek.WithPlain(plainDekBs64)
	return res, nil
}

func (i KmsManage) Encrypt(ctx context.Context, in io.Reader, out io.Writer, opts ...KmsOption) error {
	option := newKmsOptions(opts...)
	dek, err := i.cache.DekPlainGet(ctx, option.Id)
	if err != nil {
		return err
	}
	dekBs, err := base64.StdEncoding.DecodeString(dek.Plain)
	if err != nil {
		return err
	}
	// plainDek DEK明文用户缓存在内存中使用，对数据进行本地加密
	err = option.EncryptStreamFunc(in, out, dekBs, option.AESOption...)
	if err != nil {
		return err
	}
	return nil
}

func (i KmsManage) Decrypt(ctx context.Context, in io.Reader, out io.Writer, opts ...KmsOption) error {
	option := newKmsOptions(opts...)
	dek, err := i.cache.DekPlainGet(ctx, option.Id)
	if err != nil {
		return err
	}
	dekBs, err := base64.StdEncoding.DecodeString(dek.Plain)
	if err != nil {
		return err
	}
	// plainDek DEK明文用户缓存在内存中使用，对数据进行本地加密
	err = option.DecryptStreamFunc(in, out, dekBs, option.AESOption...)
	if err != nil {
		return err
	}
	return nil
}
