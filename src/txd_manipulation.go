package main

import (
	"os"
	"image"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/InfinityTools/go-squish"
)

type cachedImages struct {
	RGB, DXT1, DXT3		map[string]*[]uint8
}

var cache cachedImages

// write prepared image
func (h *txdTexture) write(f *os.File, image []uint8) error  {
	_, err := f.WriteAt(image, int64(h.Data._DataStart))
	return err
}

// prepare image to write
func (h *txdTexture) Replace(f *os.File, image image.Image) error {
	switch h.Data.TextureFormat {
	case "DXT1":
		{
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.DXT1[textureSize] != nil {
				h.write(f, *cache.DXT1[textureSize])
				//fmt.Printf("DXT1: used cache %s\n", textureSize)
			} else {
				//if cache.RGB[textureSize] != nil {
				//	squished := squish.CompressImage(cache.RGB[textureSize], squish.FLAGS_DXT1 | squish.FLAGS_ITERATIVE_CLUSTER_FIT, squish.METRIC_PERCEPTUAL)
				//	cache.DXT1[textureSize] = &squished
				//	h.write(f, *cache.DXT1[textureSize])
				//} else {
					resizedImage := imaging.Resize(image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
					cache.RGB[textureSize] = &resizedImage.Pix
					squished := squish.CompressImage(resizedImage, squish.FLAGS_DXT1 | squish.FLAGS_ITERATIVE_CLUSTER_FIT, squish.METRIC_PERCEPTUAL)
					cache.DXT1[textureSize] = &squished
					h.write(f, *cache.DXT1[textureSize])
				//}
			}
		}

	case "DXT3":
		{
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.DXT3[textureSize] != nil {
				h.write(f, *cache.DXT3[textureSize])
			} else {
				resizedImage := imaging.Resize(image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
				cache.RGB[textureSize] = &resizedImage.Pix
				squished := squish.CompressImage(resizedImage, squish.FLAGS_DXT3 | squish.FLAGS_ITERATIVE_CLUSTER_FIT, squish.METRIC_PERCEPTUAL)
				cache.DXT3[textureSize] = &squished
				h.write(f, *cache.DXT3[textureSize])
			}
		}

	case string([]byte{0x16, 0, 0, 0}), string([]byte{0x15, 0, 0, 0}):
		{
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.RGB[textureSize] != nil {
				h.write(f, *cache.RGB[textureSize])
			} else {
				resizedImage := imaging.Resize(image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
				resizedImageBGR := RGBAtoBGRA(resizedImage.Pix)
				cache.RGB[textureSize] = &resizedImageBGR
				h.write(f, *cache.RGB[textureSize])
			}
		}

	default:
		return fmt.Errorf("unknown format")
	}
	return nil
}

func RGBAtoBGRA(source []uint8) (dest []uint8) {
	length := len(source)
	dest = make([]uint8, length)
	for i := 0; i <= length - 4; i += 4 {
		dest[i]		= source[i + 2]
		dest[i + 1]	= source[i + 1]
		dest[i + 2]	= source[i]
		dest[i + 3]	= source[i + 3]
	}
	return
}

func (c *cachedImages) make(img *image.Image) {
	cache.RGB = make(map[string]*[]uint8)
	cache.DXT1 = make(map[string]*[]uint8)
	cache.DXT3 = make(map[string]*[]uint8)
	fmt.Println("Creating cache...")
	for i := 4; i <= 2048; i *= 2 {
		for j := 4; j <= 2048; j *= 2 {
			textureSize := fmt.Sprintf("%dx%d", i, j)
			resizedImage := imaging.Resize(*img, i, j, imaging.Lanczos)
			resizedImageBGR := RGBAtoBGRA(resizedImage.Pix)
			cache.RGB[textureSize] = &resizedImageBGR
			squishedDXT1 := squish.CompressImage(resizedImage, squish.FLAGS_DXT1 | squish.FLAGS_RANGE_FIT, squish.METRIC_PERCEPTUAL)
			cache.DXT1[textureSize] = &squishedDXT1
			squishedDXT3 := squish.CompressImage(resizedImage, squish.FLAGS_DXT3 | squish.FLAGS_RANGE_FIT, squish.METRIC_PERCEPTUAL)
			cache.DXT3[textureSize] = &squishedDXT3
		}
	}
	fmt.Println("Cache created")
	return
}

func (h *txdFile) replaceAll(f *os.File, image image.Image) error {
	for _, i := range h.Textures {
		i.Replace(f, image)
	}
	return nil
}