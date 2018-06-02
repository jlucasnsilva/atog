package tabbed

import (
	"log"

	"github.com/rivo/tview"
)

func newApplication() (*tview.Application, *tview.List, *tview.TextView) {
	fileList := tview.NewList()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)

	fileList.SetBorder(false)
	textView.SetBorder(false)

	ui := tview.NewFlex().
		AddItem(fileList, 0, 1, true).
		AddItem(textView, 0, 5, false)

	app := tview.NewApplication().SetRoot(ui, true)

	return app, fileList, textView
}

func runApplication(step func() string) {
	app, listView, textView := newApplication()
	con := newController(listView, textView, app)

	go func() {
		filename := ""
		for {
			if text := step(); isFileTag(text) {
				filename = getFilename(text)
				con.createComponentsIfNeeded(filename)
			} else {
				con.write(filename, text)
			}
			app.Draw()
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
