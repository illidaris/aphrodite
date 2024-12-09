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
