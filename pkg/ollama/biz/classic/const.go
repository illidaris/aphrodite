package classic

const DEF_PROMPT_TEMPLATE = `请根据以下示例学习分类规则，对新文本进行分类。
可选类别：%v
输出格式：仅返回类别名称，不包含其他内容。

示例：
%v
请对以下文本分类：
文本：%v
分类：`
