package main
import (
    gc "code.google.com/p/goncurses"
)
//Window List is the left pane in the irc client display
type WindowList struct {
    display gc.Window
    selected int
    windows []string
}

func NewWindowList(window gc.Window) *WindowList {
    w := new(WindowList)
    w.display = window
    return w
}

func (wlist *WindowList) CreateWindow(name string) {
    wlist.windows = append(wlist.windows, name)
    wlist.selectWindow(len(wlist.windows)-1)
}

func (wlist *WindowList) updateLine(id int, fmt string) {
    wlist.display.MovePrint(id, 0, fmt, wlist.windows[id])
    wlist.display.ClearToEOL()
}

func (wlist *WindowList) selectWindow(idx int) {
    wlist.updateLine(wlist.selected, " %s ")
    wlist.selected = idx
    wlist.updateLine(wlist.selected, "[%s]") 
    wlist.display.Refresh()
}
