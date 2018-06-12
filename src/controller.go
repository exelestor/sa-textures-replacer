package main

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"image"
	"path/filepath"
	"os"
)

type Message struct {
	Name		string `json:"name"`
	Params		string `json:"params"`
	TXDPath		string `json:"txd_path"`
}

var running bool

func replacerAPIHandler(request []byte) {
	if running {
		messageError(fmt.Errorf("Program is running already"))
	} else {
		running = true
		var command Message
		err := json.Unmarshal(request, &command)
		if err != nil {
			messageError(err)
			return
		}
		switch command.Name {
		case "replaceAll":
			log.Println("replacerAPIHandler: replaceAll")
			resp, err := http.Get(command.Params)
			if err != nil {
				messageError(err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				messageError(err)
				return
			}
			img, _, err := image.Decode(bytes.NewReader(body))
			if err != nil {
				messageError(err)
				return
			}
			replace(img, command.TXDPath)
		default:
			return
		}
	}
}

func replace(image image.Image, txdPath string) error {
	w.Dispatch(func() {
		w.Eval(`
document.getElementById('processLabel').style.display = 'block';
document.getElementById('processLabel').innerHTML = 'Processing...';
`)},
	)
	cache.make(&image)
	files, err := filepath.Glob(txdPath + "\\*.txd")
	fmt.Println(txdPath)
	if err != nil {
		messageError(err)
		return err
	}
	if len(files) == 0 {
		messageError(fmt.Errorf("No TXD files in this folder"))
		return nil
	}
	filesCount := len(files)
	counter := 1

	for _, fa := range files {
		fmt.Println(fa)
		progressBarSetValue(int(float64(counter)/float64(filesCount) * 100))

		f, err := os.OpenFile(fa, os.O_RDWR, 0755)
		if err != nil {
			messageError(err)
			return err
		}
		txd := new(txdFile)
		txd.read(f)

		err = txd.replaceAll(f, image)

		f.Close()
		counter++
	}

	message("Done")
	running = false
	return nil
}
