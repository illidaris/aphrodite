package classic

import "fmt"

// Label (示例)标注数据结构（对应标注文件）
type Label struct {
	Text     string `json:"text"`
	Category string `json:"category"`
}

func (l Label) GetText() string {
	return fmt.Sprintf("文本：%s\n分类：%s\n", l.Text, l.Category)
}
