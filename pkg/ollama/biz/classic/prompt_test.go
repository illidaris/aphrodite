package classic

import (
	"fmt"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	labeledData := []Label{
		{Category: "分类1", Text: "样例1ABC"},
		{Category: "类型2", Text: "样例2XXX"},
	}
	categories := []string{"分类1", "类型2"}
	newText := "待分类文本"
	// 1. 分类要求
	prompt := "请根据以下示例学习分类规则，对新文本进行分类。\n"
	prompt += fmt.Sprintf("可选类别：%v\n", categories)
	prompt += "输出格式：仅返回类别名称，不包含其他内容。\n\n"

	// 2. 标注示例（让模型学习）
	prompt += "示例：\n"
	for _, item := range labeledData {
		prompt += fmt.Sprintf("文本：%s\n分类：%s\n", item.Text, item.Category)
	}

	// 3. 待分类的新数据
	prompt += "\n请对以下文本分类：\n"
	prompt += fmt.Sprintf("文本：%s\n分类：", newText)

	v := buildPrompt(DEF_PROMPT_TEMPLATE, labeledData, categories, newText)
	if v == prompt {
		fmt.Println("测试通过")
	} else {
		t.Error("测试失败")
		println("正确： =========================================================================")
		println(v)
		println("答案： =========================================================================")
		println(v)
	}
}
