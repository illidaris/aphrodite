package base

import "time"

type Entry struct {
	Id       int64         `json:"id"`
	Code     string        `json:"code"`
	Content  string        `json:"content"`
	Prompt   string        `json:"prompt"`
	Result   string        `json:"result"`
	Duration time.Duration `json:"duration"`
}
