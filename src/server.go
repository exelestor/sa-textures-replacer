package main

import (
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
)

func replacerAPI(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request to %s from %s", r.RequestURI, r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	fmt.Println(string(body))
	replacerAPIHandler(body)
	w.Write([]byte("ok"))
}

func main() {
	log.Println("Server started")
	http.Handle("/", http.FileServer(http.Dir("./src/public_html")))
	http.HandleFunc("/api/", replacerAPI)
	if err := http.ListenAndServe(":8881", nil); err != nil {
		panic(err)
	}
}