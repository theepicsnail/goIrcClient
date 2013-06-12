package main
import (
    gc "code.google.com/p/goncurses"
)

type WindowList struct {
    windows []*Window
    display gc.Window
    selected int
}

func NewWindowList(window gc.Window) *WindowList {
    w := new(WindowList)
    w.display = window
    return w
}

func (w *WindowList) SelectWindowById(num int) *Window {
    //Unselect the old window
    w.display.Move(w.selected,0)
    w.display.ClearToEOL()
    w.display.MovePrint(w.selected,0," %s ",w.windows[w.selected].name)

    w.selected = num
    selWin := w.windows[w.selected]

    w.display.MovePrint(w.selected, 0,"[%s]", selWin.name)
    w.display.Refresh()
    return selWin
}

func (wlist *WindowList) SelectWindowByName(name string) *Window {
    id := wlist.findOrCreateWindow(name)
    wlist.SelectWindowById(id)
    return wlist.windows[id] 
}

func (wlist *WindowList) GetWindowByName(name string) *Window {
    id := wlist.findOrCreateWindow(name)
    return wlist.windows[id]
}

func (wlist *WindowList) findOrCreateWindow(name string) int {
    for n, win := range(wlist.windows) {
        if win.name == name {
            return n
        }
    }

    nWin := NewWindow(name)
    wlist.windows = append(wlist.windows, nWin)
    return len(wlist.windows) - 1
}
