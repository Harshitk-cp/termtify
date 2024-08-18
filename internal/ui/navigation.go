package ui

import (
	"github.com/gdamore/tcell/v2"
)

func (a *App) setupNavigation() {
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			a.cycleFocus()
		case tcell.KeyEsc:
			a.Stop()
		}
		return event
	})
}

func (a *App) cycleFocus() {
	currentFocus := a.GetFocus()
	switch currentFocus {
	case a.searchBar:
		a.SetFocus(a.sidebar)
	case a.sidebar:
		a.SetFocus(a.content)
	case a.content:
		a.SetFocus(a.searchBar)
	default:
		a.SetFocus(a.searchBar)
	}
}



