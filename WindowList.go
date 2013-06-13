package main
import (
    gc "code.google.com/p/goncurses"
)
//Window List is the left pane in the irc client display
type WindowList struct {
    display gc.Window
    selected int
    windows []string
    windowEvents chan *WindowEvent
}

func NewWindowList(window gc.Window) *WindowList {
    w := new(WindowList)
    w.display = window
    w.windowEvents = make(chan *WindowEvent, 4)
    go w.windowEventConsumer()
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

func (wlist *WindowList) GetWindowEventChan() chan<- *WindowEvent {
    return wlist.windowEvents
}

func (wlist *WindowList) windowEventConsumer() {
    for event := range(wlist.windowEvents) {
        switch event.eventType {
        case WIN_EVT_CREATE:
            wlist.CreateWindow(event.window.name)
        case WIN_EVT_UPDATE:
            if event.window.name != wlist.windows[wlist.selected] {
                for id, name := range(wlist.windows) {
                    if name == event.window.name {
                        wlist.updateLine(id, "*%s*")
                    }
                }
            }
            wlist.display.Refresh()
        }
    }
}

func (wlist *WindowList) SelectedWindowName() string {
    if len(wlist.windows) == 0 {
        return ""
    }
    return wlist.windows[wlist.selected]
}
