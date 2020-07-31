package utils

import (
	"fmt"
	"github.com/foobaz/lossypng/lossypng"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// 压缩图像(支持jpg/png)
func CompressImage(filename string) error {
	suffix := strings.ToLower(filepath.Ext(filename))
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		return fmt.Errorf("[CompressImage]图片格式不支持: %s", filename)
	}
	// 默认为jpg图像
	isJpg := true
	if suffix == ".png" {
		isJpg = false
	}
	// 新文件名
	newFilename := filename + ".compress"
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("文件可能不存在, err: %v", err)
	}

	// 解析图片
	var img image.Image
	if isJpg {
		img, err = jpeg.Decode(file)
	} else {
		img, err = png.Decode(file)
	}
	if err != nil {
		return fmt.Errorf("图片解析失败, err: %v", err)
	}
	file.Close()
	// 获取文件原始尺寸
	bound := img.Bounds()
	width := bound.Dx()
	height := bound.Dy()
	// 准备开始压缩
	var compressed image.Image
	if isJpg {
		// 压缩jpg, 使用Lanczos2算法进行, 无改变尺寸压缩
		compressed = resize.Resize(uint(width), uint(height), img, resize.MitchellNetravali)
	} else {
		// 压缩png, 不改变原来的色彩, 质量为原来的20%
		compressed = lossypng.Compress(img, lossypng.NoConversion, 20)
	}

	// 创建新文件
	out, err := os.Create(newFilename)
	if err != nil {
		return fmt.Errorf("创建临时文件失败, err: %v", err)
	}
	defer out.Close()

	// 编码图片
	if isJpg {
		err = jpeg.Encode(out, compressed, &jpeg.Options{Quality: 40})
	} else {
		err = png.Encode(out, compressed)
	}
	if err != nil {
		return fmt.Errorf("压缩写入失败, err: %v", err)
	}
	// 移动新文件到旧文件
	err = os.Rename(newFilename, filename)
	if err != nil {
		return fmt.Errorf("文件重命名失败, err: %v", err)
	}
	return nil
}