package main

import (
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"github.com/zserge/webview"
	"net/url"
)

var filePosition	uint32
var debug			bool
var w				webview.WebView

func main()  {
	html, err := Asset("../bin/data/index.html")
	if err != nil {
		panic(err)
	}
	css, err := Asset("../bin/data/style.css")
	if err != nil {
		panic(err)
	}
	js, err := Asset("../bin/data/script.js")
	if err != nil {
		panic(err)
	}

	w = webview.New(webview.Settings{
		Title: "sa-textures-replacer",
		URL: `data:text/html,` + url.PathEscape(string(html)),
		ExternalInvokeCallback: handleRPC,
		Width: 768,
	})
	defer w.Exit()

	w.Dispatch(func() {
		w.Bind("controller", &Controller{})
		w.InjectCSS(string(css))
		w.Eval(string(js))
	})

	w.Run()
}