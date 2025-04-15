package qrcodes

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"os"

	qrcodeReader "github.com/tuotoo/qrcode"
)

// ParseQrCode 解析二维码原始输入，支持三种输入形式：
// 1. 网络URL地址  2. 本地文件路径  3. Base64编码数据
// 参数 raw: 待解析的原始输入字符串
// 返回值: 二维码解析结果字符串，可能返回解码错误信息
func ParseQrCode(raw string) (string, error) {
	// 优先尝试解析为URL格式
	if parsedURL, err := url.Parse(raw); err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		return ReadQRCodeByUrl(raw)
	}
	// 检查是否为本地文件路径
	if _, err := os.Stat(raw); err == nil {
		return ReadQRCodeByDisk(raw)
	}
	// 默认处理为Base64编码数据
	return ReadQRCodeByBase64(raw)
}

// ReadQRCodeByUrl 从远程URL读取二维码图片进行解码
// 参数 fileurl: 远程图片资源的完整URL地址
// 返回值: 二维码解析结果字符串，可能返回网络请求或解码错误
func ReadQRCodeByUrl(fileurl string) (string, error) {
	resp, err := http.Get(fileurl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return ReadQRCodeByReader(resp.Body)
}

// ReadQRCodeByBase64 解码Base64编码的二维码图片数据
// 参数 bs64: Base64编码的图片字符串
// 返回值: 二维码解析结果字符串，可能返回Base64解码错误或二维码解析错误
func ReadQRCodeByBase64(bs64 string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(bs64)
	if err != nil {
		return "", err
	}
	return ReadQRCodeByReader(bytes.NewReader(bs))
}

// ReadQRCodeByDisk 从本地文件系统读取二维码图片文件
// 参数 filename: 本地文件路径
// 返回值: 二维码解析结果字符串，可能返回文件读取错误或解码错误
func ReadQRCodeByDisk(filename string) (string, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return ReadQRCodeByReader(bytes.NewReader(bs))
}

// ReadQRCodeByReader 从IO读取器解码二维码图片
// 参数 reader: 实现了io.Reader接口的数据源
// 返回值: 二维码解析结果字符串，可能返回解码失败错误
// 注意: 当二维码内容为空时返回空字符串且无错误
func ReadQRCodeByReader(reader io.Reader) (string, error) {
	qrmatrix, err := qrcodeReader.Decode(reader)
	if err != nil {
		return "", err
	}
	if qrmatrix == nil {
		return "", nil
	}
	return qrmatrix.Content, nil
}
