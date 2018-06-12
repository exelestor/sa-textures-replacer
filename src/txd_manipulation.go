package main

import (
	"os"
	"image"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/InfinityTools/go-squish"
	"sync"
)

type cacheField map[string]*[]uint8

type cachedImages struct {
	RGBi				map[string]*image.NRGBA
	RGB, DXT1, DXT3		cacheField
}

var cache cachedImages

// write prepared image
func (h *txdTexture) write(f *os.File, image []uint8) error  {
	_, err := f.WriteAt(image, int64(h.Data._DataStart))
	return err
}

func (h *txdTexture) writeAt(f *os.File, image []uint8, addr int64) error  {
	_, err := f.WriteAt(image, addr)
	return err
}

func (h *txdTexture) handleDXT(format uint8, image *image.Image, f *os.File) error {
	var cachedDXT *cacheField

	switch format {
	case 1:
		cachedDXT = &cache.DXT1
	case 3:
		cachedDXT = &cache.DXT3
	default:
		return fmt.Errorf("unknown compression format")
	}

	textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)

	h.write(f, *(*cachedDXT)[textureSize])

	if h.Data.MipmapCount > 1 {
		tempWidth := h.Data.Width
		tempHeight := h.Data.Height
		for _, i := range h.Data._MipmapsStart {
			tempWidth /= 2
			tempHeight /= 2
			if tempHeight < 4 || tempWidth < 4 {
				break
			}
			textureSize = fmt.Sprintf("%dx%d", tempWidth, tempHeight)
			h.writeAt(f, *(*cachedDXT)[textureSize], int64(i))
		}
	}

	return nil
}

// prepare image to write
func (h *txdTexture) Replace(f *os.File, image *image.Image) error {
	switch h.Data.TextureFormat {
	case "DXT1":
		h.handleDXT(1, image, f)

	case "DXT3":
		h.handleDXT(3, image, f)

	case string([]byte{0x16, 0, 0, 0}), string([]byte{0x15, 0, 0, 0}):
		textureSize := fmt.Sprintf("%dx%d", h.Data.Width, h.Data.Height)
		h.write(f, *cache.RGB[textureSize])

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
	cache.RGBi = make(map[string]*image.NRGBA)
	cache.RGB = make(cacheField)
	cache.DXT1 = make(cacheField)
	cache.DXT3 = make(cacheField)
	fmt.Println("Creating cache...")

	var mutex = &sync.Mutex{}

	var wgLoopRGB sync.WaitGroup
	for i := 4; i <= 2048; i *= 2 {

		j := 4

		for j <= 2048 {
			wgLoopRGB.Add(10)

			for worker := 1; worker <= 10; worker++ {
				go func(x int, y int) {
					defer wgLoopRGB.Done()
					textureSize := fmt.Sprintf("%dx%d", x, y)
					resizedImage := imaging.Resize(*img, x, y, imaging.Lanczos)
					resizedImageBGR := RGBAtoBGRA(resizedImage.Pix)
					mutex.Lock()
					cache.RGBi[textureSize] = resizedImage
					cache.RGB[textureSize] = &resizedImageBGR
					mutex.Unlock()
					fmt.Println(textureSize, "RGB")
				}(i, j)
				j *= 2
			}

			wgLoopRGB.Wait()
		}
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()

		var wgLoop1 sync.WaitGroup
		for i := 4; i <= 2048; i *= 2 {

			j := 4

			for j <= 2048 {
				wgLoop1.Add(5)

				for worker := 1; worker <= 5; worker++ {
					go func(x int, y int) {
						defer wgLoop1.Done()
						textureSize := fmt.Sprintf("%dx%d", x, y)
						squishedDXT1 := squish.CompressImage(cache.RGBi[textureSize], squish.FLAGS_DXT1 | squish.FLAGS_RANGE_FIT | squish.FLAGS_SOURCE_BGRA, squish.METRIC_PERCEPTUAL)
						mutex.Lock()
						cache.DXT1[textureSize] = &squishedDXT1
						mutex.Unlock()
						fmt.Println(textureSize, "DXT1")
					}(i, j)
					j *= 2
				}

				wgLoop1.Wait()
			}
		}
	}()

	go func() {
		defer wg.Done()

		var wgLoop1 sync.WaitGroup
		for i := 4; i <= 2048; i *= 2 {

			j := 4

			for j <= 2048 {
				wgLoop1.Add(5)

				for worker := 1; worker <= 5; worker++ {
					go func(x int, y int) {
						defer wgLoop1.Done()
						textureSize := fmt.Sprintf("%dx%d", x, y)
						squishedDXT3 := squish.CompressImage(cache.RGBi[textureSize], squish.FLAGS_DXT3 | squish.FLAGS_RANGE_FIT | squish.FLAGS_SOURCE_BGRA, squish.METRIC_PERCEPTUAL)
						mutex.Lock()
						cache.DXT3[textureSize] = &squishedDXT3
						mutex.Unlock()
						fmt.Println(textureSize, "DXT3")
					}(i, j)
					j *= 2
				}

				wgLoop1.Wait()
			}
		}
	}()

	wg.Wait()

	fmt.Println("Cache created")
	return
}

func (h *txdFile) replaceAll(f *os.File, image image.Image) error {
	for _, i := range h.Textures {
		i.Replace(f, &image)
	}
	return nil
}