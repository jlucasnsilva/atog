package ui

import (
	"log"

	"github.com/jlucasnsilva/atog/stalker"
	"github.com/rivo/tview"
)

// Execute ...
func Execute(filenames []string, values <-chan stalker.Value) {
	fileList := tview.NewList()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)

	fileList.SetBorder(true)
	textView.SetBorder(true)

	ui := tview.NewFlex().
		AddItem(fileList, 0, 1, true).
		AddItem(textView, 0, 6, false)

	app := tview.NewApplication()

	go func() {
		router := newRouter(filenames, textView, fileList)
		for v := range values {
			router.handle(v)
			app.Draw()
		}
	}()

	if err := app.SetRoot(ui, true).Run(); err != nil {
		log.Fatalln(err)
	}
}
