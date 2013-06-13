package main

type WindowEventType int
const (
    WIN_EVT_CREATE WindowEventType = iota
    WIN_EVT_UPDATE
    WIN_EVT_FOCUS //Make this window the main window
)

type WindowEvent struct {
    window *Window
    eventType WindowEventType
}

// Window structs serve to model an individual private message/room/etc.
type Window struct {
    name string
    history []string
    lineChan chan string
    handler chan<- *WindowEvent
}

func NewWindow(name string, eventStream chan<- *WindowEvent) *Window {
    win := new(Window)
    win.name = name
    win.lineChan = make(chan string)
    win.handler = eventStream
    go win.lineReader()
    
    eventStream <- win.newEvent(WIN_EVT_CREATE)

    return win
}

func (win *Window) newEvent(evtType WindowEventType) *WindowEvent {
    return &WindowEvent{win, evtType}
}

func (win *Window) GetLineChan() chan<- string {
    return win.lineChan
}

func (win *Window) lineReader() {
    for line := range(win.lineChan) {
        win.history = append(win.history, line)

        go func() { win.handler <- win.newEvent(WIN_EVT_UPDATE) }() 
    }
}
