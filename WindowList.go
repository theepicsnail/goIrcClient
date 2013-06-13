package main
import (
    gc "code.google.com/p/goncurses"
)
//Window List is the left pane in the irc client display
type WindowList struct {
    display gc.Window
    lastSelected int
}

func NewWindowList(window gc.Window) *WindowList {
    w := new(WindowList)
    w.display = window
    return w
}

func (wlist *WindowList) Update() {
    wlist.display.Erase()
/*
    selectedIdx := wlist.manager.selected

    for n, window := range(wlist.manager.windows) {
        fmtStr := " %s "
        if n == selectedIdx {
            fmtStr = "[%s]"
        }
        wlist.display.MovePrint(n, 0, fmtStr, window.name) 
    }
*/
    wlist.display.Refresh()
}
