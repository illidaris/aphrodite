package qrcodes

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/nfnt/resize"
)

// ImageWithLogo 将底图和logo合并
func ImageWithLogo(src, logo io.Reader, logoP int, out io.Writer) error {
	if logoP > 10 || logoP < 0 {
		return fmt.Errorf("logo占比必须在0-10之间")
	}
	// 将二维码文件接码成图片
	srcImg, srcName, srcErr := image.Decode(src)
	if srcErr != nil {
		return srcErr
	}
	if srcImg == nil {
		return fmt.Errorf("source %v %v", srcName, srcErr)
	}
	// 将填充图解码成png图片
	logoImg, logoName, logoErr := image.Decode(logo)
	if logoErr != nil {
		return logoErr
	}
	if logoImg == nil {
		return fmt.Errorf("logo %v %v", logoName, srcErr)
	}
	// 调整Logo大小（二维码大小的1/5）
	logoSize := srcImg.Bounds().Dx() / logoP
	resizedLogo := resize.Resize(uint(logoSize), 0, logoImg, resize.Lanczos3)
	// 计算Logo位置（居中）
	offset := image.Pt(
		(srcImg.Bounds().Dx()-resizedLogo.Bounds().Dx())/2,
		(srcImg.Bounds().Dy()-resizedLogo.Bounds().Dy())/2,
	)
	// 创建画布并合并图片
	canvas := image.NewRGBA(srcImg.Bounds())
	draw.Draw(canvas, srcImg.Bounds(), srcImg, image.Point{}, draw.Over)
	draw.Draw(canvas, resizedLogo.Bounds().Add(offset), resizedLogo, image.Point{}, draw.Over)
	return png.Encode(out, canvas)
}

/*
ReadFile 根据输入类型自动选择读取方式
参数:
  - raw: 输入字符串，可以是URL/文件路径/Base64编码数据

返回值:
  - []byte: 读取到的文件内容字节数组
  - error: 读取过程中发生的错误（包含各种处理方式的错误传递）
*/
func ReadFile(raw string) ([]byte, error) {
	// 尝试解析为HTTP/HTTPS URL格式
	if parsedURL, err := url.Parse(raw); err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		return ReadFileByUrl(raw)
	}

	// 检查是否存在对应的本地文件
	if _, err := os.Stat(raw); err == nil {
		return ReadFileByDisk(raw)
	}

	// 默认处理为Base64编码数据
	return ReadFileByBase64(raw)
}

/*
ReadFileByUrl 通过HTTP GET请求获取远程文件内容
参数:
  - fileurl: 完整的文件URL地址（需包含http/https协议头）

返回值:
  - []byte: 下载的文件内容字节数组
  - error: 网络请求失败或响应异常时返回错误
*/
func ReadFileByUrl(fileurl string) ([]byte, error) {
	resp, err := http.Get(fileurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

/*
ReadFileByBase64 解码Base64格式字符串
参数:
  - bs64: 符合RFC 4648标准的Base64编码字符串

返回值:
  - []byte: 解码后的原始字节数据
  - error: 输入不符合Base64格式时返回解码错误
*/
func ReadFileByBase64(bs64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(bs64)
}

/*
ReadFileByDisk 读取本地文件系统上的文件
参数:
  - filename: 本地文件的绝对或相对路径

返回值:
  - []byte: 文件内容字节数组
  - error: 文件不存在或读取权限不足时返回错误
*/
func ReadFileByDisk(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
