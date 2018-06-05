package main

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"image"
)

type Message struct {
	Name		string `json:"name"`
	Params		string `json:"params"`
}

func replacerAPIHandler(request []byte) {
	var command Message
	err := json.Unmarshal(request, &command)
	if err != nil {
		log.Println("error json unmarshal")
		return
	}
	switch command.Name {
	case "replaceAll":
		//replace()
		log.Println("replacerAPIHandler: replaceAll")
	case "downloadImage":
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
	default:
		return
	}
}
