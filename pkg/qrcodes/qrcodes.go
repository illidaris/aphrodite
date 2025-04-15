package qrcodes

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"

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
