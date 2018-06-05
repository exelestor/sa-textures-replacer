package main

import (
	"fmt"
	"os"
	"flag"
	"log"
	"path/filepath"

	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"image"

	"github.com/disintegration/imaging"
)

var filePosition uint32
var debug bool
var onlyRead *bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func replace() error {
	//files, err := filepath.Glob("F:\\txd\\*.txd")
	files, err := filepath.Glob("bin/test.txd")
	check(err)
	filesCount := len(files)
	counter := 1

	for _, fa := range files {
		fmt.Printf("[%d/%d] Working with '%s'... ", counter, filesCount, fa)
		f, err := os.OpenFile(fa, os.O_RDWR, 0755)
		check(err)
		txd := new(txdFile)
		txd.read(f)

		if !(*onlyRead) {
			err := txd.replaceAll(f)
			if err != nil {
				fmt.Println("Some errors", err)
			} else {
				fmt.Println("Done")
			}
		} else {
			fmt.Println("Done")
		}
		f.Close()
		counter++
	}

	return nil
}

/*
	TODO: error handling
 */

func not_main() {
	debugFlag := flag.Bool("debug", false, "a bool")
	onlyRead = flag.Bool("read", false, "a bool")
	flag.Parse()
	debug = *debugFlag

	fmt.Println("Шалом")

	reader, err := os.Open("bin/test.jpg")
	if err != nil {
	    log.Fatal(err)
	}
	defer reader.Close()

	n, _, err := image.Decode(reader)
	//onePixel := n.At(1, 1)
	//r, g, b, a := onePixel.RGBA()
	//
	//fmt.Println(r, g, b, a)
	//fmt.Printf("%x %x %x %x\n", r >> 8, g >> 8, b >> 8, a >> 8)
	//fmt.Printf("%b %b %b %b\n", r >> 8, g >> 8, b >> 8, a >> 8)

	imagesBuffer := make(map[string]*image.NRGBA)

	minResolution := 4
	maxResolution := 4096

	//uiprogress.Start()            // start rendering
	//bar := uiprogress.AddBar(11 * 11) // Add a new bar
	//
	//bar.AppendCompleted()
	//bar.PrependElapsed()

	for i := minResolution; i <= maxResolution; i = i << 1 {
		for j := minResolution; j <= maxResolution; j = j << 1 {
			fmt.Print(i, j)
			imagesBuffer[fmt.Sprintf("%dx%d", i, j)] = imaging.Resize(n, i, j, imaging.Lanczos)
			fmt.Println(" done")
		}
	}
	pixels := (*imagesBuffer["16x16"]).Pix

	fmt.Printf("%+v\n", pixels[:4])

	//replace()

	fmt.Println("All done!")
}
