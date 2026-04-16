package services

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"
)

type ImageService interface {
	SaveImage(file io.Reader, filename string, size int64) (string, error)
	DeleteImage(path string) error
}

type imageService struct{}

func NewImageService() ImageService {
	return &imageService{}
}

func (s *imageService) SaveImage(file io.Reader, filename string, size int64) (string, error) {

	// ensure directory exists
	err := os.MkdirAll("static/uploads", os.ModePerm)
	if err != nil {
		return "", err
	}

	// sanitize filename (basic)
	ext := filepath.Ext(filename)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	path := filepath.Join("static/uploads", name)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 🔥 compress if large
	if size > 600*1024 {
		img, format, err := image.Decode(file)
		if err != nil {
			return "", fmt.Errorf("invalid image")
		}

		switch format {
		case "jpeg", "jpg":
			err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 70})
		case "png":
			encoder := png.Encoder{
				CompressionLevel: png.BestCompression,
			}
			err = encoder.Encode(dst, img)
		default:
			return "", fmt.Errorf("unsupported image type")
		}

		if err != nil {
			return "", err
		}

	} else {
		_, err = io.Copy(dst, file)
		if err != nil {
			return "", err
		}
	}

	// return public path
	return "/uploads/" + name, nil
}

func (s *imageService) DeleteImage(path string) error {

	if path == "" {
		return nil // nothing to delete
	}

	// 🔥 Safety check — prevent accidental deletion outside uploads
	if !filepath.HasPrefix(path, "/uploads/") {
		return fmt.Errorf("invalid image path")
	}

	// convert public path → filesystem path
	fullPath := filepath.Join("static", path)

	err := os.Remove(fullPath)
	if err != nil {
		// ignore if file already doesn't exist
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}