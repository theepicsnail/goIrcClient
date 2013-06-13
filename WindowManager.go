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
    if num < 0 || num >= len(w.windows) {
        return nil
    }

    w.selected = num
    win := w.windows[w.selected]

    w.eventChannel <- &WindowEvent{win, WIN_EVT_FOCUS}

    return win
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
