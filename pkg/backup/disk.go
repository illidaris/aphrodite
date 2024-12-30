package backup

import (
	"context"
	"encoding/json"
	"os"
)

func DiskSave(ctx context.Context, fullname string, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(fullname, bs, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func DiskLoad(ctx context.Context, fullname string, data interface{}) error {
	bs, err := os.ReadFile(fullname)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, data)
}

func ReadFrmDisk(fullname string) string {
	bs, err := os.ReadFile(fullname)
	if err != nil {
		return ""
	}
	return string(bs)
}

func WriteToDisk(fullname string, content string) {
	_ = os.WriteFile(fullname, []byte(content), os.ModePerm)
}
