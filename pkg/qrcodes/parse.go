package qrcodes

import (
	"bytes"
	"errors"
	"image"
	"io"

	"github.com/makiuchi-d/gozxing"
	gozxingQrCode "github.com/makiuchi-d/gozxing/qrcode"
)

// ParseQrCode 解析二维码原始输入，支持三种输入形式：
// 1. 网络URL地址  2. 本地文件路径  3. Base64编码数据
// 参数 raw: 待解析的原始输入字符串
// 返回值: 二维码解析结果字符串，可能返回解码错误信息
func ParseQrCode(raw string) (string, error) {
	bs, err := ReadFile(raw)
	if err != nil {
		return "", err
	}
	return ReadQRCodeByReader(bytes.NewBuffer(bs))
}

// ReadQRCodeByReader 从IO读取器解码二维码图片
// 参数 reader: 实现了io.Reader接口的数据源
// 返回值: 二维码解析结果字符串，可能返回解码失败错误
// 注意: 当二维码内容为空时返回空字符串且无错误
func ReadQRCodeByReader(reader io.Reader) (string, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}
	// image转bmp
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	// 解码
	qrReader := gozxingQrCode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil) // Safe for concurrent use
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("decode result is nil")
	}
	return result.GetText(), nil
}
