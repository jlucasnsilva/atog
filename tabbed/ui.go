package tabbed

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func newApplication() (*tview.Application, *tview.List, *tview.TextView) {
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

	app := tview.NewApplication().SetRoot(ui, true)
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyCtrlL {
			textView.Clear()
			app.Draw()
			return nil
		}
		return e
	})

	return nil, fileList, textView
}

func runApplication(step func() string) {
	app, listView, textView := newApplication()
	con := newController(listView, textView)

	go func() {
		for {
			text := step()
			if isFileTag(text) {
				con.switchTo(getFilename(text))
			} else {
				con.write(text)
			}
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
