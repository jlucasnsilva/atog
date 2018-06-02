package atog

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// NewTextViewApp ...
func NewTextViewApp() (*tview.Application, *tview.TextView) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)

	app := tview.NewApplication()
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyCtrlL {
			textView.Clear()
			app.Draw()
			return nil
		}
		return e
	})
	return app, textView
}

// RunSimpleApp ...
func RunSimpleApp(step func() string) {
	app, textView := NewTextViewApp()

	go func() {
		fmt.Fprintln(textView, Highlight(step()))
		app.Draw()
	}()

	if err := app.SetRoot(textView, true).Run(); err != nil {
		log.Fatalln(err)
	}
}

// RunSimpleLoopApp ...
func RunSimpleLoopApp(step func() string) {
	app, textView := NewTextViewApp()

	go func() {
		for {
			fmt.Fprintln(textView, step())
			app.Draw()
		}
	}()

	if err := app.SetRoot(textView, true).Run(); err != nil {
		log.Fatalln(err)
	}
}
