package main

type WindowEvent struct {
    window *Window
    line string
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
    return win
}

func (win *Window) GetLineChan() chan<- string {
    return win.lineChan
}

func (win *Window) lineReader() {
    for line := range(win.lineChan) {
        win.history = append(win.history, line)

        event := new(WindowEvent)
        event.window = win
        event.line = line
        go func() { win.handler <- event }() 
    }
}
