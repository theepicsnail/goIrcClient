package main
import (
    gc "code.google.com/p/goncurses"
  //  "fmt"
)
type WindowList struct {
    windowNames []string
    display gc.Window
    selected int
}
func NewWindowList(window gc.Window) *WindowList {
    w := new(WindowList)
    w.display = window
    w.createWindow("win1")
    w.createWindow("win2")
    return w
}

func (w *WindowList) createWindow(name string) {
    w.windowNames = append(w.windowNames, name)
    w.selectWindowById(len(w.windowNames)-1)
}

func (w *WindowList) selectWindowById(num int) {
    w.display.Move(w.selected,0)
    w.display.ClearToEOL()
    w.display.MovePrint(w.selected,0," %s ",w.windowNames[w.selected])
    w.selected = num
    w.display.MovePrint(w.selected, 0,"[%s]", w.windowNames[w.selected])
    w.display.Refresh()
}
