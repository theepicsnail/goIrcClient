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
    wlist.SelectWindowById(len(wlist.windows)-1)
}

func (wlist *WindowList) updateLine(id int, fmt string) {
    go func() {CURSES.Lock()
    wlist.display.MovePrint(id, 0, fmt, wlist.windows[id])
    wlist.display.ClearToEOL()
    wlist.display.Refresh()
    CURSES.Unlock()}()
}

func (wlist *WindowList) SelectWindowById(idx int) {
    wlist.updateLine(wlist.selected, " %s ")
    wlist.selected = idx
    wlist.updateLine(wlist.selected, "[%s]") 
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
            if event.window.name != wlist.SelectedWindowName() {
                id := wlist.findIdByName(event.window.name)
                wlist.updateLine(id, "*%s*")
            }
        case WIN_EVT_FOCUS:
            wlist.SelectWindowById(wlist.findIdByName(event.window.name))
        }
    }
}

func (wlist *WindowList) findIdByName(name string) int {
    for id, wname := range(wlist.windows) {
        if wname == name {
            return id
        }
    }
    return -1
}

func (wlist *WindowList) SelectedWindowName() string {
    if len(wlist.windows) == 0 {
        return ""
    }
    return wlist.windows[wlist.selected]
}
