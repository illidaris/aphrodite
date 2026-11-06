package gormex

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gorm.io/gorm/schema"
)

func TestEncryptedSerializerValue(t *testing.T) {
	ctx := context.Background()
	serializer := EncryptedSerializer{EncryptFunc: func(_ context.Context, value string) (string, error) {
		return "enc:" + value, nil
	}}

	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{name: "nil remains nil", input: nil, want: nil},
		{name: "empty string remains empty", input: "", want: ""},
		{name: "string is encrypted", input: "secret", want: "enc:secret"},
		{name: "other type is ignored", input: 42, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializer.Value(ctx, nil, reflect.Value{}, tt.input)
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestEncryptedSerializerValueWithoutEncryptFunc(t *testing.T) {
	_, err := (EncryptedSerializer{}).Value(context.Background(), nil, reflect.Value{}, "secret")
	if err == nil {
		t.Fatal("Value() error = nil, want missing EncryptFunc error")
	}
}

func TestEncryptedSerializerScan(t *testing.T) {
	ctx := context.Background()
	var got string
	field := &schema.Field{Set: func(_ context.Context, dst reflect.Value, value interface{}) error {
		dst.Elem().SetString(value.(string))
		return nil
	}}
	serializer := EncryptedSerializer{DecryptFunc: func(_ context.Context, value string) (string, error) {
		return "plain:" + value, nil
	}}
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{name: "nil leaves field unchanged", input: nil, want: "initial"},
		{name: "empty string clears field", input: "", want: ""},
		{name: "string is decrypted", input: "cipher", want: "plain:cipher"},
		{name: "bytes are decrypted", input: []byte("cipher"), want: "plain:cipher"},
		{name: "unsupported type leaves field unchanged", input: 42, want: "initial"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got = "initial"
			if err := serializer.Scan(ctx, field, reflect.ValueOf(&got), tt.input); err != nil {
				t.Fatalf("Scan() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Scan() set %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEncryptedSerializerScanErrors(t *testing.T) {
	ctx := context.Background()
	field := &schema.Field{Set: func(_ context.Context, dst reflect.Value, value interface{}) error {
		dst.Elem().SetString(value.(string))
		return nil
	}}
	var dst string
	if err := (EncryptedSerializer{}).Scan(ctx, field, reflect.ValueOf(&dst), "cipher"); err == nil {
		t.Fatal("Scan() error = nil, want missing DecryptFunc error")
	}
	wantErr := errors.New("decrypt failed")
	serializer := EncryptedSerializer{DecryptFunc: func(context.Context, string) (string, error) {
		return "", wantErr
	}}
	if err := serializer.Scan(ctx, field, reflect.ValueOf(&dst), "cipher"); !errors.Is(err, wantErr) {
		t.Errorf("Scan() error = %v, want %v", err, wantErr)
	}
}
