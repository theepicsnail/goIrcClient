package main

// Window structs serve to model an individual private message/room/etc.
type Window struct {
    name string
    history []string
}

func NewWindow(name string) *Window {
    win := new(Window)
    win.name = name
    return win
}

func (win *Window) AddLine(line string) {
    win.history = append(win.history, line)
}
