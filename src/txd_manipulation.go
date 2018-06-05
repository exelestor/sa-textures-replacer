package main

import (
	"os"
	"image"
	"fmt"
	"github.com/disintegration/imaging"
)

type cachedImages struct {
	RGB, DXT1, DXT3		map[string]*[]uint8
}

var cache cachedImages

func convertToDXT1(image []uint8) (ret []uint8) {

	return
}

func convertToDXT3(image []uint8) (ret []uint8) {

	return
}

// write prepared image
func (h *txdTexture) write(f *os.File, image []uint8) error  {
	_, err := f.WriteAt(image, int64(h.Data._DataStart))
	return err
}

// prepare image to write
func (h *txdTexture) Replace(f *os.File, image *image.Image) error {
	switch h.Data.TextureFormat {
	case string(0x21):
		{
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.RGB[textureSize] != nil {
				h.write(f, *cache.RGB[textureSize])
			} else {
				resizedImage := imaging.Resize(*image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
				*cache.RGB[textureSize] = resizedImage.Pix
				h.write(f, *cache.RGB[textureSize])
			}
		}

	case "DXT1":
		{
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.DXT1[textureSize] != nil {
				h.write(f, *cache.DXT1[textureSize])
			} else {
				if cache.RGB[textureSize] != nil {
					*cache.DXT1[textureSize] = convertToDXT1(*cache.RGB[textureSize])
					h.write(f, *cache.DXT1[textureSize])
				} else {
					resizedImage := imaging.Resize(*image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
					*cache.RGB[textureSize] = resizedImage.Pix
					*cache.DXT1[textureSize] = convertToDXT1(*cache.RGB[textureSize])
					h.write(f, *cache.DXT1[textureSize])
				}
			}
		}

	case "DXT3":
		{
			textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
			if cache.DXT3[textureSize] != nil {
				h.write(f, *cache.DXT3[textureSize])
			} else {
				if cache.RGB[textureSize] != nil {
					*cache.DXT3[textureSize] = convertToDXT3(*cache.RGB[textureSize])
					h.write(f, *cache.DXT3[textureSize])
				} else {
					resizedImage := imaging.Resize(*image, int(h.Data.Width), int(h.Data.Height), imaging.Lanczos)
					*cache.RGB[textureSize] = resizedImage.Pix
					*cache.DXT3[textureSize] = convertToDXT3(*cache.RGB[textureSize])
					h.write(f, *cache.DXT3[textureSize])
				}
			}
		}

	default:
		return nil
	}

	return nil
}

func (h *txdFile) replaceAll(f *os.File) error {
	for _, _ = range h.Textures {
		//i.Replace(f, )
	}
	return nil
}