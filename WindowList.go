package main
import (
    gc "code.google.com/p/goncurses"
)

type Window struct {
    name string
    history []string
}

type WindowList struct {
    windows []Window
    display gc.Window
    selected int
}

func NewWindowList(window gc.Window) *WindowList {
    w := new(WindowList)
    w.display = window
    return w
}

func (w *WindowList) SelectWindowById(num int) Window {
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

func (wlist *WindowList) SelectWindowByName(name string) Window {
    //Finds a window by name, if it's not found, create it.
    id := wlist.findWindowByName(name)

    if id == -1 {
        id = wlist.createWindow(name)
    }

    return wlist.SelectWindowById(id)
}

func (wList *WindowList) findWindowByName(name string) int {
    for n, win := range(wList.windows) {
        if win.name == name {
            return n
        }
    }
    return -1
}

func (wlist *WindowList) createWindow(name string) int {
    win := new(Window)
    win.name = name
    win.history = make([]string, 16)
    wlist.windows = append(wlist.windows, *win)
    return len(wlist.windows) - 1
}
