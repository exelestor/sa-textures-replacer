package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"
	_ "image/gif"
	_ "image/jpeg"
	"log"
	"image/png"
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
	files, err := filepath.Glob("test.txd")
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
				fmt.Println("Some errors")
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
func main() {
	debugFlag := flag.Bool("debug", false, "a bool")
	onlyRead = flag.Bool("read", false, "a bool")
	flag.Parse()
	debug = *debugFlag

	reader, err := os.Open("test.png")
	if err != nil {
	    log.Fatal(err)
	}
	defer reader.Close()

	n, err := png.Decode(reader)
	fmt.Printf("%+v\n", n)
	fmt.Println(n.At(1, 1))

	for i := 4; i <= 4096; i = i << 1 {
		for j := 4; j <= 4096; j = j << 1 {
			fmt.Printf("%dx%d\n", i, j)
		}
	}

	replace()

	fmt.Println("All done!")
}
