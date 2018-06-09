package main

import (
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"github.com/zserge/webview"
	"net/url"
	"fmt"
	"strings"
	"log"
)

var filePosition uint32
var debug bool

type Controller struct {
	TXDDir		string `json:"txdDir"`
	Progress	string `json:"progress"`
}

func check(e error) {
	if e != nil {

		w.Dispatch(func() {
			w.Eval(fmt.Sprintf(
				`document.getElementById('processLabel').innerHTML = '<span class="pink">Error:</span> %s'`,
				e.Error(),
			))
			w.Eval(`document.getElementById('processLabel').style.display = 'block';`)
		},
		)
		log.Println(e)
	}
}

func (c *Controller) Replace(txdPath string, picUrl string) {
	fmt.Println(txdPath, picUrl)
	if txdPath == "" {
		check(fmt.Errorf("Specify TXD directory"))
	} else if picUrl == "" {
		check(fmt.Errorf("Specify Picture URL"))
	} else {
		replaceJSON := fmt.Sprintf(
			`{"name":"replaceAll","txd_path":"%s","params":"%s"}`,
			strings.Replace(txdPath, "\\", "\\\\", -1),
			picUrl,
		)

		go replacerAPIHandler([]byte(replaceJSON))
	}
}

func handleRPC(w webview.WebView, data string) {
	switch data {
	case "opendir":
		path := w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "Open directory", "")
		w.Eval(fmt.Sprintf(`document.getElementById('txdDir').value = '%s';`, strings.Replace(path, "\\", "\\\\", -1)))
	}
}

var w webview.WebView

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