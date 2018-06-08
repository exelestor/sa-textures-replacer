package main

import (
	"fmt"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"github.com/andlabs/ui"
)

var filePosition uint32
var debug bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main()  {
	debug = false
	err := ui.Main(func() {
		txdpath		:= ui.NewEntry()
		picurl		:= ui.NewEntry()
		submit		:= ui.NewButton("Replace")
		progressBar	:= ui.NewProgressBar()
		greeting	:= ui.NewLabel("")
		box			:= ui.NewVerticalBox()
		box.Append(ui.NewLabel("txd path:"), false)
		box.Append(txdpath, false)
		box.Append(ui.NewLabel("pic url:"), false)
		box.Append(picurl, false)
		box.Append(submit, false)
		box.Append(progressBar, false)
		window := ui.NewWindow("sa-textures-replacer", 200, 100, true)
		window.SetMargined(true)
		window.SetChild(box)
		submit.OnClicked(func(*ui.Button) {
			greeting.SetText("Working...")
			progressBar.SetValue(0)
			fmt.Printf("%v\n", progressBar)
			go replacerAPIHandler([]byte(fmt.Sprintf(`{"name":"replaceAll","params":"%s","txd_path":"%s"}`, picurl.Text(), txdpath.Text())), progressBar)
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}