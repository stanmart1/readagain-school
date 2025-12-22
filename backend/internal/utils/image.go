package utils

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageOptimizer struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
}

func NewImageOptimizer(maxWidth, maxHeight, quality int) *ImageOptimizer {
	return &ImageOptimizer{
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
		Quality:   quality,
	}
}

func (io *ImageOptimizer) OptimizeImage(filePath string) error {
	img, err := imaging.Open(filePath)
	if err != nil {
		return NewInternalServerError("Failed to open image", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width > io.MaxWidth || height > io.MaxHeight {
		img = imaging.Fit(img, io.MaxWidth, io.MaxHeight, imaging.Lanczos)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	
	switch ext {
	case ".jpg", ".jpeg":
		return io.saveJPEG(img, filePath)
	case ".png":
		return io.saveJPEG(img, strings.TrimSuffix(filePath, ext)+".jpg")
	default:
		return io.saveJPEG(img, filePath)
	}
}

func (io *ImageOptimizer) saveJPEG(img image.Image, filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return NewInternalServerError("Failed to create optimized image", err)
	}
	defer out.Close()

	opts := &jpeg.Options{Quality: io.Quality}
	if err := jpeg.Encode(out, img, opts); err != nil {
		return NewInternalServerError("Failed to encode image", err)
	}

	return nil
}
