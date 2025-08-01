package classic

import "fmt"

// 生成提示词（包含标注示例和分类要求）
func buildPrompt(template string, labeledData []Label, categories []string, newText string) string {
	// 2. 标注示例（让模型学习）
	labelTxt := ""
	for _, item := range labeledData {
		labelTxt += item.GetText()
	}
	if template == "" {
		template = DEF_PROMPT_TEMPLATE
	}
	return fmt.Sprintf(template, categories, labelTxt, newText)
}
