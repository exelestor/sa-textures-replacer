package main

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"image"
	"github.com/andlabs/ui"
	"path/filepath"
	"os"
)

type Message struct {
	Name		string `json:"name"`
	Params		string `json:"params"`
	TXDPath		string `json:"txd_path"`
}

func replacerAPIHandler(request []byte, bar *ui.ProgressBar) {
	var command Message
	err := json.Unmarshal(request, &command)
	if err != nil {
		log.Println("error json unmarshal")
		return
	}
	switch command.Name {
	case "replaceAll":
		log.Println("replacerAPIHandler: replaceAll")
		resp, err := http.Get(command.Params)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		img, typeImg, err := image.Decode(bytes.NewReader(body))
		fmt.Printf("Size:	%d\n", len(body))
		fmt.Printf("Bounds:	%s\n", img.Bounds())
		fmt.Printf("Type:	%s\n", typeImg)
		fmt.Printf("%v\n", bar)
		replace(img, command.TXDPath, bar)
	default:
		return
	}
}

func replace(image image.Image, txdPath string, bar *ui.ProgressBar) error {
	go cache.make(&image)
	files, err := filepath.Glob(txdPath + "\\*.txd")
	fmt.Println(txdPath)
	check(err)
	filesCount := len(files)
	counter := 1

	for _, fa := range files {
		//fmt.Printf("[%d/%d] Working with '%s'... ", counter, filesCount, fa)
		ui.QueueMain(func() {
			bar.SetValue(int(float64(counter) / float64(filesCount) * 100))
		})
		f, err := os.OpenFile(fa, os.O_RDWR, 0755)
		check(err)
		txd := new(txdFile)
		txd.read(f)

		err = txd.replaceAll(f, image)
		//if err != nil {
		//	fmt.Println("Some errors", err)
		//} else {
		//	fmt.Println("Done")
		//}

		f.Close()
		counter++
	}

	return nil
}
