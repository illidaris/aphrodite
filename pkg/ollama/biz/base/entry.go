package base

type Entry struct {
	Id       int64  `json:"id"`
	Code     string `json:"code"`
	Content  string `json:"content"`
	Prompt   string `json:"prompt"`
	Result   string `json:"result"`
	Duration int64  `json:"duration"`
}
