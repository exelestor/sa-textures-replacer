package main

import (
	"os"
	"image"
	"fmt"
	"github.com/disintegration/imaging"
	"log"
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
	case string(0x22), string(0x21):
		{
			log.Println("32bit")
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.RGB[textureSize] != nil {
				h.write(f, *cache.RGB[textureSize])
			} else {
				resizedImage := imaging.Resize(image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
				cache.RGB[textureSize] = &resizedImage.Pix
				h.write(f, *cache.RGB[textureSize])
			}
		}

	case "DXT1":
		{
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.DXT1[textureSize] != nil {
				h.write(f, *cache.DXT1[textureSize])
			} else {
				//if cache.RGB[textureSize] != nil {
				//	squished := squish.CompressImage(cache.RGB[textureSize], squish.FLAGS_DXT1 | squish.FLAGS_ITERATIVE_CLUSTER_FIT, squish.METRIC_PERCEPTUAL)
				//	cache.DXT1[textureSize] = &squished
				//	h.write(f, *cache.DXT1[textureSize])
				//} else {
					resizedImage := imaging.Resize(image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
					cache.RGB[textureSize] = &resizedImage.Pix
					squished := squish.CompressImage(resizedImage, squish.FLAGS_DXT1 | squish.FLAGS_RANGE_FIT, squish.METRIC_PERCEPTUAL)
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
				//if cache.RGB[textureSize] != nil {
				//	squished := squish.CompressImage(image, squish.FLAGS_DXT3 | squish.FLAGS_ITERATIVE_CLUSTER_FIT, squish.METRIC_PERCEPTUAL)
				//	cache.DXT3[textureSize] = &squished
				//	h.write(f, *cache.DXT3[textureSize])
				//} else {
					resizedImage := imaging.Resize(image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
					cache.RGB[textureSize] = &resizedImage.Pix
					squished := squish.CompressImage(resizedImage, squish.FLAGS_DXT3 | squish.FLAGS_RANGE_FIT, squish.METRIC_PERCEPTUAL)
					cache.DXT3[textureSize] = &squished
					h.write(f, *cache.DXT3[textureSize])
				//}
			}
		}

	default:
		{
			log.Println("default")
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
	}

	return nil
}

func RGBAtoBGRA(image []uint8) (ret []uint8) {
	length := len(image)
	ret = make([]uint8, length)
	for i := 0; i <= length - 4; i += 4 {
		ret[i]		= image[i + 2]
		ret[i + 1]	= image[i + 1]
		ret[i + 2]	= image[i]
		ret[i + 3]	= image[i + 3]
	}
	return
}

func (h *txdFile) replaceAll(f *os.File, image image.Image) error {
	cache.RGB = make(map[string]*[]uint8)
	cache.DXT1 = make(map[string]*[]uint8)
	cache.DXT3 = make(map[string]*[]uint8)
	for _, i := range h.Textures {
		i.Replace(f, image)
	}
	return nil
}