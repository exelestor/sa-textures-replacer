package main

import (
	"fmt"
	"github.com/zserge/webview"
	"log"
	"strings"
)

type Controller struct {
	TXDDir   string `json:"txdDir"`
	Progress string `json:"progress"`
}

func messageError(e error) {
	if e != nil {

		w.Dispatch(
			func() {
				w.Eval(
					fmt.Sprintf(
						`document.getElementById('processLabel').innerHTML = '<span class="pink">Error:</span> %s'`,
						e.Error(),
					))
				w.Eval(`document.getElementById('processLabel').style.display = 'block';`)
			},
		)
		log.Println(e)
	}
}

func message(s string) {
	w.Dispatch(
		func() {
			w.Eval(
				fmt.Sprintf(
					`document.getElementById('processLabel').innerHTML = '%s'`, s,
				),
			)
			w.Eval(`document.getElementById('processLabel').style.display = 'block';`)
		},
	)
}

func progressBarSetValue(p int) {
	w.Dispatch(
		func() {
			w.Eval(
				fmt.Sprintf("document.getElementById('myBar').style.width = %d + '%%';", p),
			)
		},
	)
}

func (c *Controller) Replace(txdPath string, picUrl string) {
	fmt.Println(txdPath, picUrl)
	if txdPath == "" {
		messageError(fmt.Errorf("Specify TXD directory"))
	} else if picUrl == "" {
		messageError(fmt.Errorf("Specify Picture URL"))
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
