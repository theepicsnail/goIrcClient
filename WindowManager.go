package main

type WindowManager struct {
    windows []*Window
    selected int
    eventChannel chan<- *WindowEvent
}

func NewWindowManager(evtChan chan<- *WindowEvent) *WindowManager {
    wm := new(WindowManager)
    wm.eventChannel = evtChan
    return wm
}

func (w *WindowManager) SelectWindowById(num int) *Window {
    w.selected = num
    return w.windows[w.selected]
}

func (wlist *WindowManager) SelectWindowByName(name string) *Window {
    id := wlist.findOrCreateWindow(name)
    return wlist.SelectWindowById(id)
}

func (wlist *WindowManager) GetWindowByName(name string) *Window {
    id := wlist.findOrCreateWindow(name)
    return wlist.windows[id]
}

func (wlist *WindowManager) findOrCreateWindow(name string) int {
    for n, win := range(wlist.windows) {
        if win.name == name {
            return n
        }
    }

    nWin := NewWindow(name, wlist.eventChannel)
    wlist.windows = append(wlist.windows, nWin)
    return len(wlist.windows) - 1
}
