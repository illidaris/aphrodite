package qrcodes

import (
	qrcode "github.com/skip2/go-qrcode"
)

// WriteQrCode 生成指定内容的二维码并输出到目标文件或返回字节数组
// 参数:
//   - content: 需要编码的二维码内容字符串
//   - quality: 二维码纠错等级(使用qrcode包定义的RecoveryLevel类型)
//   - size: 生成二维码的像素尺寸(正方形边长)
//   - dest: 目标文件路径，为空时返回字节数组，非空时写入文件
//
// 返回值:
//   - []byte: 当dest为空时返回PNG格式字节数据
//   - error: 生成过程中可能发生的错误
func WriteQrCode(content string, quality qrcode.RecoveryLevel, size int, dest string) ([]byte, error) {
	// 当指定目标路径时，直接写入文件并返回nil字节数组
	if dest != "" {
		return nil, qrcode.WriteFile(content, qrcode.Medium, size, dest)
	}

	// 未指定路径时，生成PNG格式字节数组
	bs, err := qrcode.Encode(content, quality, size)
	if err != nil {
		return bs, err
	}
	return bs, nil
}
