package encrypter

import (
	"encoding/json"
	"testing"
)

// TestStringEncrypt 测试 StringEncrypt 函数
func TestStringEncrypt(t *testing.T) {
	opts := []DataEncryptOptionFunc{
		WithSecretString("0123456789abcdef"),
	}
	raw := "test data sccc"
	enStr, err := StringEncrypt(raw, opts...)
	if err != nil {
		t.Errorf("EncryptStream failed: %v", err)
	}
	println("RawEn: " + enStr)

	deStr, deErr := StringDecrypt(enStr, opts...)
	if err != nil {
		t.Errorf("DecryptStream failed: %v", deErr)
	}
	println("RawDe: " + deStr)

	if deStr != raw {
		t.Error("DecryptStream failed")
	}
}

// TestJsonEncrypt 测试 JsonEncrypt 函数
func TestJsonEncrypt(t *testing.T) {
	opts := []DataEncryptOptionFunc{
		WithSecretString("asdda@#￥@12阿"),
	}
	type Sub struct {
		Data1 string `json:"data1"`
		Data2 string `json:"data2"`
	}
	type Demo struct {
		Name   string  `json:"name"`
		Age    int     `json:"age"`
		Amount float64 `json:"amount"`
		Data   Sub     `json:"data"`
	}
	demo := Demo{
		Name:   "张三",
		Age:    18,
		Amount: 123.123,
		Data: Sub{
			Data1: "data1",
			Data2: "data2",
		},
	}
	enStr, err := JsonEncrypt(demo, opts...)
	if err != nil {
		t.Errorf("EncryptStream failed: %v", err)
	}
	println("RawEn: " + enStr)

	deV, deErr := JsonDecrypt[Demo](enStr, opts...)
	if err != nil {
		t.Errorf("DecryptStream failed: %v", deErr)
	}
	deStrBs, _ := json.Marshal(deV)
	println("RawDe: " + string(deStrBs))

	if deV != demo {
		t.Error("DecryptStream failed")
	}
}
